package main

import (
	"flag"
	"fmt"

	"github.com/xogroup/kapacitor-configmap-listener/configuration"
	"github.com/xogroup/kapacitor-configmap-listener/helpers"
)

func main() {
	kubeConfig := flag.String("kubeconfig", "", "absolute path to the kubeconfig file (optional) [\"~/.kube/config\"]")
	inCluster := flag.Bool("incluster", false, "setup context for inside cluster (optional) [false]")
	flag.Parse()

	// creates the clientset
	kubeClient, err := configuration.NewClientSet(inCluster, kubeConfig)
	if err != nil {
		panic(err.Error())
	}

	helpers.Watch(
		kubeClient,
		"services",
		func(obj interface{}) {
			fmt.Printf("service added: %s \n", obj)
		},
		func(obj interface{}) {
			fmt.Printf("service deleted: %s \n", obj)
		},
		func(oldObj, newObj interface{}) {
			fmt.Printf("service changed \n")
		})
}
