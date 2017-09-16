package kapacitor

import (
	"fmt"

	"github.com/xogroup/kapacitor-configmap-listener/templates"

	"github.com/influxdata/kapacitor/client/v1"
	"k8s.io/client-go/pkg/api/v1"

	log "github.com/sirupsen/logrus"
)

// TaskEntry stores the context of the scaling ConfigMap from k8s.  It also records if the
// entry has been processed.
type TaskEntry struct {
	vars map[string]client.Var
}

type taskWork struct {
	taskOptions client.CreateTaskOptions
	action      ActionType
}

// TaskStore all of the individual Task in a map which can be looked up via the ReleaseName
// of a deployment
type TaskStore struct {
	kapacitorClient *client.Client
	templateID      string
	Store           map[string]*TaskEntry
	workQueue       chan taskWork
}

// ActionType signals what behavior is desired
type ActionType int

const (
	Create ActionType = iota
	Update
	Delete
)

// NewTaskStore instantiates a new object of that type
func NewTaskStore(kapacitorClient *client.Client) (*TaskStore, error) {

	defaultOptions := &client.ListTasksOptions{
		Limit: 500,
	}

	tasks, err := kapacitorClient.ListTasks(defaultOptions)
	if err != nil {
		return nil, err
	}

	store := map[string]*TaskEntry{}

	for key := range tasks {
		store[tasks[key].ID] = &TaskEntry{
			vars: tasks[key].Vars,
		}
	}

	log.Infof("Found %d task in Kapacitor (%s)\n", len(store), kapacitorClient.URL())

	return &TaskStore{
		kapacitorClient: kapacitorClient,
		Store:           store,
		workQueue:       make(chan taskWork),
	}, nil
}

// CreateTask converts the configMap into a Kapacitor task and adds it to the
// worker queue to be added by the processor
func (taskStore *TaskStore) CreateTask(configMap *v1.ConfigMap) error {

	log.Infoln("Creating Task")
	return taskStore.pushTask(configMap, Create)
}

// UpdateTask converts the configMap into a Kapacitor task and adds it to the
// worker queue to be updated by the processor
func (taskStore *TaskStore) UpdateTask(configMap *v1.ConfigMap) error {

	log.Infoln("Updating Task")
	return taskStore.pushTask(configMap, Update)
}

// DeleteTask converts the configMap into a Kapacitor task and adds it to the
// worker queue to be removed by the processor
func (taskStore *TaskStore) DeleteTask(configMap *v1.ConfigMap) error {

	log.Infoln("Deleting Task")
	return taskStore.pushTask(configMap, Delete)
}

func (taskStore *TaskStore) pushTask(configMap *v1.ConfigMap, action ActionType) error {

	id := configMap.Data["releaseName"]

	taskOptions, err := buildTaskOptions(configMap)
	if err != nil {
		return err
	}

	taskStore.Store[id] = &TaskEntry{
		vars: taskOptions.Vars,
	}

	go func() {
		taskStore.workQueue <- taskWork{
			taskOptions: *taskOptions,
			action:      action,
		}
	}()

	return nil
}

func buildTaskOptions(configMap *v1.ConfigMap) (*client.CreateTaskOptions, error) {

	if template, ok := tick.Templates[configMap.Data["template"]]; ok {
		return &client.CreateTaskOptions{
			ID: configMap.Data["releaseName"],
			// this should be moved into a factory func()
			TemplateID: template.ID,
			Vars:       buildVars(configMap),
		}, nil
	}

	return nil, fmt.Errorf("no TICK template found with name of %s", configMap.Data["template"])
}

func buildVars(configMap *v1.ConfigMap) map[string]client.Var {

	vars := map[string]client.Var{}

	vars["database"] = client.Var{Type: client.VarString, Value: configMap.Data["database"]}
	vars["retentionPolicy"] = client.Var{Type: client.VarString, Value: configMap.Data["retentionPolicy"]}
	vars["measurement"] = client.Var{Type: client.VarString, Value: configMap.Data["measurement"]}
	//where_filter
	vars["field"] = client.Var{Type: client.VarString, Value: configMap.Data["field"]}
	vars["target"] = client.Var{Type: client.VarFloat, Value: configMap.Data["target"]}
	vars["deploymentName"] = client.Var{Type: client.VarFloat, Value: configMap.Data["deploymentName"]}
	vars["scalingCooldown"] = client.Var{Type: client.VarDuration, Value: configMap.Data["scalingCooldown"]}
	vars["descalingCooldown"] = client.Var{Type: client.VarDuration, Value: configMap.Data["descalingCooldown"]}

	return vars
}
