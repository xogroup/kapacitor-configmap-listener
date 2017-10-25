package handlers

import (
	"github.com/xogroup/kapacitor-configmap-listener/helpers/influx"
	"github.com/xogroup/kapacitor-configmap-listener/helpers/kapacitor"

	log "github.com/sirupsen/logrus"
)

// KapacitorResetHandler is an object to hold shared context for handling Kapacitor resets
type KapacitorResetHandler struct {
	taskStore     *kapacitor.TaskStore
	subscriptions *influx.Subscriptions
}

// NewKapacitorResetHandlers instantiates a new object of that type
func NewKapacitorResetHandlers(taskStore *kapacitor.TaskStore, subscriptions *influx.Subscriptions) *KapacitorResetHandler {
	return &KapacitorResetHandler{taskStore, subscriptions}
}

// Handle comparing kapacitor state and resetting it when necessary
func (context *KapacitorResetHandler) Handle() {

	log.Debugln("Kapacitor Sync Polled")

	isSync, err := context.taskStore.IsSync()

	if err != nil {
		log.Errorf("Error syncing Kapacitor (%v)", err)
		return
	}

	if !isSync {
		log.Infoln("Kapacitor and TaskStore not Sync!!!")
		if context.subscriptions != nil {
			context.subscriptions.RemoveAll()
		}

		context.taskStore.Reseed()
	} else {
		log.Debugln("Kapacitor and TaskStore in Sync")
	}
}
