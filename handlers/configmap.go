package handlers

import (
	"strings"

	"github.com/xogroup/kapacitor-configmap-listener/helpers/kapacitor"

	"k8s.io/client-go/pkg/api/v1"

	log "github.com/sirupsen/logrus"
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

	log.Infof("ConfigMap (%s) created", obj.(*v1.ConfigMap).Name)
	filterAndProcess(obj, context.prefix, context.taskStore.CreateTask)
}

//HandleUpdated captures created config map events and re-processes the rule to Kapacitor
func (context *ConfigMapHandlers) HandleUpdated(oldObj interface{}, newObj interface{}) {

	log.Infof("ConfigMap (%s) updated", newObj.(*v1.ConfigMap).Name)
	filterAndProcess(newObj, context.prefix, context.taskStore.UpdateTask)
}

//HandleDeleted captures created config map events and deletes the rule from Kapacitor
func (context *ConfigMapHandlers) HandleDeleted(obj interface{}) {

	log.Infof("ConfigMap (%s) deleted", obj.(*v1.ConfigMap).Name)
	filterAndProcess(obj, context.prefix, context.taskStore.DeleteTask)
}

func filterAndProcess(obj interface{}, prefix string, f func(*v1.ConfigMap) error) {
	configMap := obj.(*v1.ConfigMap)

	if strings.HasPrefix(configMap.ObjectMeta.Name, prefix) {
		err := f(configMap)
		if err != nil {
			log.WithError(err)
		}
	}
}
