package main

import (
	"flag"
	"fmt"

	"k8s.io/client-go/pkg/api/v1"

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
		"configmaps",
		func(obj interface{}) {
			configMap := obj.(*v1.ConfigMap)
			fmt.Printf("configmap created: %s \n", configMap.ObjectMeta.Name)
		},
		func(obj interface{}) {
			configMap := obj.(*v1.ConfigMap)
			fmt.Printf("configmap deleted: %s \n", configMap.ObjectMeta.Name)
		},
		func(oldObj, newObj interface{}) {
			configMap := oldObj.(*v1.ConfigMap)
			fmt.Printf("configmap deleted: %s \n", configMap.ObjectMeta.Name)
		})
}
