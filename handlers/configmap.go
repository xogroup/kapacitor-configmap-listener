package handlers

import (
	"fmt"
	"strings"

	"github.com/xogroup/kapacitor-configmap-listener/helpers/kapacitor"

	"k8s.io/client-go/pkg/api/v1"
)

// ConfigMapHandlers is an object to hold shared context for all of the handlers
type ConfigMapHandlers struct {
	prefix    string
	taskStore *kapacitor.TaskStore
}

//NewConfigMapHandlers instantiates a new object of that type
func NewConfigMapHandlers(prefix string, taskStore *kapacitor.TaskStore) *ConfigMapHandlers {
	return &ConfigMapHandlers{prefix, taskStore}
}

//HandleCreated captures created config map events and processes it as new rules to Kapacitor
func (context *ConfigMapHandlers) HandleCreated(obj interface{}) {

	filterAndProcess(obj, context.prefix, context.taskStore.CreateTask)

	// fmt.Printf("configmap created: %s \n", configMap.ObjectMeta.Name)
	// fmt.Println("------------------------------------------------------")
	// fmt.Println(configMap.Data)
	// fmt.Println("======================================================")
}

//HandleUpdated captures created config map events and re-processes the rule to Kapacitor
func (context *ConfigMapHandlers) HandleUpdated(oldObj interface{}, newObj interface{}) {

	filterAndProcess(newObj, context.prefix, context.taskStore.UpdateTask)
	// fmt.Printf("configmap updated: %s \n", configMap.ObjectMeta.Name)
	// fmt.Println("------------------------------------------------------")
	// fmt.Println(configMap.Data)
	// fmt.Println("======================================================")
}

//HandleDeleted captures created config map events and deletes the rule from Kapacitor
func (context *ConfigMapHandlers) HandleDeleted(obj interface{}) {

	filterAndProcess(obj, context.prefix, context.taskStore.DeleteTask)
	// fmt.Printf("configmap deleted: %s \n", configMap.ObjectMeta.Name)
	// fmt.Println("------------------------------------------------------")
	// fmt.Println(configMap.Data)
	// fmt.Println("======================================================")
}

func filterAndProcess(obj interface{}, prefix string, f func(*v1.ConfigMap) error) {
	configMap := obj.(*v1.ConfigMap)

	if strings.HasPrefix(configMap.ObjectMeta.Name, prefix) {
		err := f(configMap)

		if err != nil {
			fmt.Println(err.Error())
		}
	}
}
