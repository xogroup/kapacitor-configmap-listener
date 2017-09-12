package helpers

import (
	"time"

	"github.com/xogroup/kapacitor-configmap-listener/factories"

	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/tools/cache"
)

// Watch sets up the watch list against the configmap resource and uses the provided event
// handlers for events when bubbled.
func Watch(kubeClient *kubernetes.Clientset, resource string, addHandler func(obj interface{}), deleteHandler func(obj interface{}), updateHandler func(oldObj, newObj interface{})) chan struct{} {

	watchlist := cache.NewListWatchFromClient(kubeClient.Core().RESTClient(), resource, v1.NamespaceAll, fields.Everything())

	_, controller := cache.NewInformer(
		watchlist,
		factories.NewSpecFactory().Build(resource),
		time.Second*0,
		cache.ResourceEventHandlerFuncs{
			AddFunc:    addHandler,
			DeleteFunc: deleteHandler,
			UpdateFunc: updateHandler,
		},
	)

	stop := make(chan struct{})

	go controller.Run(stop)

	for {
		time.Sleep(time.Second)
	}
}
