package main

import (
	"flag"
	"os"

	"github.com/xogroup/kapacitor-configmap-listener/configuration"
	"github.com/xogroup/kapacitor-configmap-listener/handlers"
	"github.com/xogroup/kapacitor-configmap-listener/helpers/k8s"
	"github.com/xogroup/kapacitor-configmap-listener/helpers/kapacitor"

	log "github.com/sirupsen/logrus"
)

func main() {
	kubeConfig := flag.String("kubeconfig", "", "absolute path to the kubeconfig file (optional) [\"~/.kube/config\"]")
	inCluster := flag.Bool("incluster", false, "setup context for inside cluster (optional) [false]")
	prefix := flag.String("prefixname", "kapacitor-hpa-rule-", "prefix name to capture for event handling for config maps (optional) [\"kapacitor-hpa-rule-\"]")
	kapacitorURL := flag.String("kapacitorurl", os.Getenv("KAPACITOR_URL"), "url path to the kapacitord server.  Defaults to the KAPACITOR_URL environment variable if set (optional) [\"localhost:9092\"]")
	logLevel := flag.Int("loglevel", 4, "log level 0-5 {panic, fatal, error warn, info, debug} (optional) [4-info]")

	flag.Parse()

	log.SetLevel(log.Level(uint32(*logLevel)))

	// creates the clientset
	kubeClient, err := configuration.NewKubeClient(inCluster, kubeConfig)
	if err != nil {
		panic(err.Error())
	}

	// create kapacitor client
	kapacitorClient, err := configuration.NewKapacitorClient(*kapacitorURL)
	if err != nil {
		panic(err.Error())
	}

	// create local storage for desired and real state
	taskStore, err := kapacitor.NewTaskStore(kapacitorClient)
	if err != nil {
		panic(err.Error())
	}

	//check to see if kapacitor is up
	//list kapacitor tasks and keep it in memory
	// name of tasks should be release names

	configMapHandlers := handlers.NewConfigMapHandlers(*prefix, taskStore)

	k8s.Watch(
		kubeClient,
		"configmaps",
		configMapHandlers.HandleCreated,
		configMapHandlers.HandleDeleted,
		configMapHandlers.HandleUpdated)
}
