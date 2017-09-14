package handlers

import (
	"fmt"
	"strings"

	"k8s.io/client-go/pkg/api/v1"
)

// ConfigMapHandlers is an object to hold shared context for all of the handlers
type ConfigMapHandlers struct {
	prefix string
}

//NewConfigMapHandlers instantiates a new object of that type
func NewConfigMapHandlers(prefix string) *ConfigMapHandlers {
	return &ConfigMapHandlers{prefix}
}

//HandleCreated captures created config map events and processes it as new rules to Kapacitor
func (context *ConfigMapHandlers) HandleCreated(obj interface{}) {
	configMap := obj.(*v1.ConfigMap)
	if strings.HasPrefix(configMap.ObjectMeta.Name, context.prefix) {

		fmt.Println(configMap.Data["target"])
	}

	fmt.Printf("configmap created: %s \n", configMap.ObjectMeta.Name)
	fmt.Println("------------------------------------------------------")
	fmt.Println(configMap.Data)
	fmt.Println("======================================================")
}

//HandleUpdated captures created config map events and re-processes the rule to Kapacitor
func (context *ConfigMapHandlers) HandleUpdated(oldObj interface{}, newObj interface{}) {
	configMap := oldObj.(*v1.ConfigMap)
	fmt.Printf("configmap updated: %s \n", configMap.ObjectMeta.Name)
	fmt.Println("------------------------------------------------------")
	fmt.Println(configMap.Data)
	fmt.Println("======================================================")
}

//HandleDeleted captures created config map events and deletes the rule from Kapacitor
func (context *ConfigMapHandlers) HandleDeleted(obj interface{}) {
	configMap := obj.(*v1.ConfigMap)
	fmt.Printf("configmap deleted: %s \n", configMap.ObjectMeta.Name)
	fmt.Println("------------------------------------------------------")
	fmt.Println(configMap.Data)
	fmt.Println("======================================================")
}
