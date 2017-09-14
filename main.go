package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/xogroup/kapacitor-configmap-listener/configuration"
	"github.com/xogroup/kapacitor-configmap-listener/handlers"
	"github.com/xogroup/kapacitor-configmap-listener/helpers"
)

func main() {
	kubeConfig := flag.String("kubeconfig", "", "absolute path to the kubeconfig file (optional) [\"~/.kube/config\"]")
	inCluster := flag.Bool("incluster", false, "setup context for inside cluster (optional) [false]")
	prefix := flag.String("prefixname", "kapacitor-hpa-rule-", "prefix name to capture for event handling for config maps (optional) [\"kapacitor-hpa-rule-\"]")
	kapacitorURL := flag.String("kapacitorurl", os.Getenv("KAPACITOR_URL"), "url path to the kapacitord server.  Defaults to the KAPACITOR_URL environment variable if set (optional) [\"localhost:9092\"]")
	flag.Parse()

	// creates the clientset
	kubeClient, err := configuration.NewKubeClient(inCluster, kubeConfig)
	if err != nil {
		panic(err.Error())
	}

	// create kapacitor client
	kapacitorClient, err := configuration.NewKapacitorClient(*kapacitorURL)

	_, text, err := kapacitorClient.Ping()

	if err != nil {
		panic(err.Error())
	}

	fmt.Println(text)

	//check to see if kapacitor is up
	//list kapacitor tasks and keep it in memory
	// name of tasks should be release names

	configMapHandlers := handlers.NewConfigMapHandlers(*prefix)

	helpers.Watch(
		kubeClient,
		"configmaps",
		configMapHandlers.HandleCreated,
		configMapHandlers.HandleDeleted,
		configMapHandlers.HandleUpdated)
}
