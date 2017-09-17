package kapacitor

import (
	"bytes"
	"fmt"

	"github.com/xogroup/kapacitor-configmap-listener/templates"

	"github.com/influxdata/kapacitor/client/v1"
	"k8s.io/client-go/pkg/api/v1"

	log "github.com/sirupsen/logrus"
)

// TaskEntry stores the context of the scaling ConfigMap from k8s.  It also records if the
// entry has been processed.
type TaskEntry struct {
	name      string
	namespace string
	vars      client.Vars
}

type work struct {
	taskOptions *TaskOptions
	taskEntry   *TaskEntry
	action      ActionType
}

// TaskStore all of the individual Task in a map which can be looked up via the ReleaseName
// of a deployment
type TaskStore struct {
	kapacitorClient *client.Client
	templateID      string
	Store           map[string]*TaskEntry
	workQueue       chan work
}

// ActionType signals what behavior is desired
type ActionType int

const (
	Create ActionType = iota
	Update
	Delete
)

// TaskOptions is a generic store for kapacitor.client.CreateTaskOptions and kapacitor.client.UpdateTaskOptions
type TaskOptions struct {
	ID         string
	TemplateID string
	DBRPs      []client.DBRP
	Vars       client.Vars
	Status     client.TaskStatus
	TICKscript string
	Type       client.TaskType
}

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

	taskStore := &TaskStore{
		kapacitorClient: kapacitorClient,
		Store:           store,
		workQueue:       make(chan work),
	}

	go taskStore.workProcessor()

	return taskStore, nil
}

func (taskStore *TaskStore) workProcessor() {

	kapacitorClient := taskStore.kapacitorClient

	for {
		select {
		case job := <-taskStore.workQueue:
			go func() {
				preExistingTaskLink := kapacitorClient.TaskLink(job.taskOptions.ID)
				preExistingTask, _ := kapacitorClient.Task(preExistingTaskLink, nil)

				switch job.action {
				case Create, Update:

					if preExistingTask.ID == "" {
						task, err := kapacitorClient.CreateTask(*job.taskOptions.ToCreateTaskOptions())
						if err != nil {
							log.Errorf("Task %s (%v)", job.taskOptions.ID, err)
							return
						}

						log.Infof("Task %s created with status of %s", task.ID, task.Status)
					} else {
						task, err := kapacitorClient.UpdateTask(preExistingTaskLink, *job.taskOptions.ToUpdateTaskOptions())
						if err != nil {
							log.Errorf("Task %s (%v)", job.taskOptions.ID, err)
							return
						}

						log.Infof("Task %s updated with status of %s", task.ID, task.Status)
					}

				case Delete:

					if preExistingTask.ID != "" {
						err := kapacitorClient.DeleteTask(preExistingTaskLink)
						if err != nil {
							log.Errorf("Task %s (%v)", job.taskOptions.ID, err)
							return
						}

						log.Infof("Task %s deleted")
					} else {
						log.Infof("Task %s does not exist in kapacitor")
					}
				}
				log.Infof("Processed job for task %s", job.taskOptions.ID)
			}()
		}
	}
}

// CreateTask converts the configMap into a Kapacitor task and adds it to the
// worker queue to be added by the processor
func (taskStore *TaskStore) CreateTask(configMap *v1.ConfigMap) error {

	log.Infof("Creating task %s.%s", configMap.Namespace, configMap.Name)
	return taskStore.pushTask(configMap, Create)
}

// UpdateTask converts the configMap into a Kapacitor task and adds it to the
// worker queue to be updated by the processor
func (taskStore *TaskStore) UpdateTask(configMap *v1.ConfigMap) error {

	log.Infof("Updating task %s.%s", configMap.Namespace, configMap.Name)
	return taskStore.pushTask(configMap, Update)
}

// DeleteTask converts the configMap into a Kapacitor task and adds it to the
// worker queue to be removed by the processor
func (taskStore *TaskStore) DeleteTask(configMap *v1.ConfigMap) error {

	log.Infof("Deleting task %s.%s", configMap.Namespace, configMap.Name)
	return taskStore.pushTask(configMap, Delete)
}

func (taskStore *TaskStore) pushTask(configMap *v1.ConfigMap, action ActionType) error {

	id := configMap.Data["releaseName"]
	taskOptions, err := buildTaskOptions(configMap)
	if err != nil {
		return err
	}

	taskEntry := &TaskEntry{
		name:      configMap.Name,
		namespace: configMap.Namespace,
		vars:      taskOptions.Vars,
	}

	taskStore.Store[id] = taskEntry

	go func() {
		taskStore.workQueue <- work{
			taskOptions: taskOptions,
			taskEntry:   taskEntry,
			action:      action,
		}
	}()

	return nil
}

func buildTaskOptions(configMap *v1.ConfigMap) (*TaskOptions, error) {

	vars, err := buildVars(configMap)

	if err != nil {
		return nil, err
	}

	templateID := vars["template"].Value.(string)

	if template, ok := tick.Templates[templateID]; ok {

		dbrp := buildDBRP(vars)

		return &TaskOptions{
			ID: vars["releaseName"].Value.(string),
			// TemplateID: template.ID,
			DBRPs:      *dbrp,
			Vars:       vars,
			Status:     client.Disabled,
			TICKscript: template.Template,
			Type:       client.StreamTask,
		}, nil
	}

	return nil, fmt.Errorf("no TICK template found with name of %s", templateID)
}

func buildDBRP(vars client.Vars) *[]client.DBRP {

	dbrps := []client.DBRP{
		client.DBRP{
			Database:        vars["database"].Value.(string),
			RetentionPolicy: vars["retentionPolicy"].Value.(string),
		},
	}

	return &dbrps
}

func buildVars(configMap *v1.ConfigMap) (client.Vars, error) {

	var jsonBuffer bytes.Buffer
	vars := client.Vars{}
	index := 1
	length := len(configMap.Data)

	jsonBuffer.WriteString("{")

	for key := range configMap.Data {
		jsonBuffer.WriteString(fmt.Sprintf("\"%s\":%s", key, configMap.Data[key]))

		if index < length {
			jsonBuffer.WriteString(",")
		}

		index++
	}

	jsonBuffer.WriteString("}")

	err := vars.UnmarshalJSON(jsonBuffer.Bytes())

	return vars, err
}

// ToUpdateTaskOptions converts the generic TaskOptions to the kapacitor.client.UpdateTaskOptions specific type
func (taskOptions *TaskOptions) ToUpdateTaskOptions() *client.UpdateTaskOptions {
	return &client.UpdateTaskOptions{
		ID:         taskOptions.ID,
		TemplateID: taskOptions.TemplateID,
		DBRPs:      taskOptions.DBRPs,
		Vars:       taskOptions.Vars,
		Status:     taskOptions.Status,
		TICKscript: taskOptions.TICKscript,
		Type:       taskOptions.Type,
	}
}

// ToCreateTaskOptions converts the generic TaskOptions to the kapacitor.client.CreateTaskOptions specific type
func (taskOptions *TaskOptions) ToCreateTaskOptions() *client.CreateTaskOptions {
	return &client.CreateTaskOptions{
		ID:         taskOptions.ID,
		TemplateID: taskOptions.TemplateID,
		DBRPs:      taskOptions.DBRPs,
		Vars:       taskOptions.Vars,
		Status:     taskOptions.Status,
		TICKscript: taskOptions.TICKscript,
		Type:       taskOptions.Type,
	}
}
