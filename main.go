package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/xogroup/kapacitor-configmap-listener/configuration"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
	for {
		pods, err := kubeClient.CoreV1().Pods("").List(metav1.ListOptions{})
		if err != nil {
			panic(err.Error())
		}
		fmt.Printf("There are %d pods in the cluster\n", len(pods.Items))
		time.Sleep(10 * time.Second)
	}
}
