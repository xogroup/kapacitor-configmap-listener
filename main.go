package main

import (
	"flag"
	"os"
	"strconv"
	"time"

	"github.com/xogroup/kapacitor-configmap-listener/configuration"
	"github.com/xogroup/kapacitor-configmap-listener/handlers"
	"github.com/xogroup/kapacitor-configmap-listener/helpers/influx"
	"github.com/xogroup/kapacitor-configmap-listener/helpers/k8s"
	"github.com/xogroup/kapacitor-configmap-listener/helpers/kapacitor"

	log "github.com/sirupsen/logrus"
)

func main() {
	kubeConfig := flag.String("kubeconfig", "", "absolute path to the kubeconfig file (optional) [\"~/.kube/config\"]")
	inCluster := flag.Bool("incluster", false, "setup context for inside cluster (optional) [false]")
	prefix := flag.String("prefixname", "kapacitor-hpa-rule-", "prefix name to capture for event handling for config maps (optional) [\"kapacitor-hpa-rule-\"]")
	kapacitorURL := flag.String("kapacitorurl", os.Getenv("KAPACITOR_URL"), "url path to the kapacitord server.  Defaults to KAPACITOR_URL environment variable if set (optional) [\"localhost:9092\"]")
	logLevel := flag.Int("loglevel", 4, "log level 0-5 {panic, fatal, error warn, info, debug} (optional) [4-info]")
	cleanupSubscriptions := flag.Bool("cleansubscriptions", false, "delete all subscriptions from influx cluster (optional) [false]")
	influxURL := flag.String("influxurl", os.Getenv("INFLUX_URL"), "url path to the influxdb server.  Defaults to INFLUX_URL environment variable if set (optional) [\"localhost:8086\"]")
	influxUserName := flag.String("influxusername", os.Getenv("INFLUX_USERNAME"), "username for influx.  Defaults to INFLUX_USERNAME environment variable if set (optional)")
	influxPassword := flag.String("influxpassword", os.Getenv("INFLUX_PASSWORD"), "password for influx.  Defaults to INFLUX_PASSWORD environment variable if set (optional)")
	influxSsl := flag.Bool("influxssl", stringToBool(os.Getenv("INFLUX_SSL")), "use ssl for influx. Defaults to INFLUX_SSL environment variable if set (optional) [false]")
	influxUnsafeSsl := flag.Bool("influxunsafessl", stringToBool(os.Getenv("INFLUX_UNSAFE_SSL")), "skip ssl verification. Defaults to INFLUX_UNSAFE_SSL environment variable if set (optional) [false]")

	var subscriptions *influx.Subscriptions

	flag.Parse()

	log.SetLevel(log.Level(uint32(*logLevel)))

	if *cleanupSubscriptions == true {
		influxClient, err := configuration.NewInfluxClient(*influxURL, *influxUserName, *influxPassword, *influxSsl, *influxUnsafeSsl)
		if err != nil {
			panic(err.Error())
		}

		subscriptions = influx.NewSubscriptions(influxClient)
	}

	// creates kubernetes client
	kubeClient, err := configuration.NewKubeClient(inCluster, kubeConfig)
	if err != nil {
		panic(err.Error())
	}

	// create kapacitor client
	kapacitorClient, err := configuration.NewKapacitorClient(*kapacitorURL)
	if err != nil {
		panic(err.Error())
	}

	// create local storage for kapacitor state alignment
	taskStore, err := kapacitor.NewTaskStore(kapacitorClient)
	if err != nil {
		panic(err.Error())
	}

	// initialize config map handlers for valid change events
	configMapHandlers := handlers.NewConfigMapHandlers(*prefix, taskStore)

	// create a watch for config map changes using the event handlers
	k8s.Watch(
		kubeClient,
		"configmaps",
		configMapHandlers.HandleCreated,
		configMapHandlers.HandleDeleted,
		configMapHandlers.HandleUpdated)

	// intitialize kapacitor reset handler for polling events
	kapacitorResetHandler := handlers.NewKapacitorResetHandlers(taskStore, subscriptions)

	// create a watcher that polls kapacitor for changes
	kapacitor.Watch(kapacitorResetHandler.Handle)

	for {
		time.Sleep(time.Second)
	}
}

func stringToBool(value string) bool {
	result, err := strconv.ParseBool(value)

	if err != nil {
		return false
	}

	return result
}
