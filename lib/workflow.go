package spiderswarm

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

type ActionTemplate struct {
	Name              string
	StructName        string
	ConstructorParams map[string]interface{} // XXX: should this be defined more strictly?
}

type DataPipeTemplate struct {
	SourceActionName string
	SourceOutputName string
	DestActionName   string
	DestInputName    string
	TaskInputName    string
	TaskOutputName   string
}

type TaskTemplate struct {
	TaskName          string
	Initial           bool
	ActionTemplates   []ActionTemplate
	DataPipeTemplates []DataPipeTemplate
}

type Workflow struct {
	Name          string
	Version       string
	TaskTemplates []TaskTemplate
}

func (w *Workflow) Run() ([]*Item, error) {
	jobUUID := uuid.New().String()
	startedAt := time.Now()

	var items []*Item

	log.Info(fmt.Sprintf("Job %s started from workflow %s:%s at %v", jobUUID, w.Name, w.Version,
		startedAt))

	var tasks []*Task
	var task *Task

	for _, taskTempl := range w.TaskTemplates {
		if !taskTempl.Initial {
			continue
		}

		newTask := NewTaskFromTemplate(&taskTempl, w, jobUUID) // TODO: implement this one

		tasks = append(tasks, newTask)
	}

	for {
		if len(tasks) == 0 {
			break
		}

		task, tasks = tasks[0], tasks[1:]
		log.Info(fmt.Sprintf("Running task %v", task))
		err := task.Run()
		if err != nil {
			log.Error(fmt.Sprintf("Task %v failed with error: %v", task, err))
		} else { // TODO: make this less nested
			for _, outDP := range task.Outputs {
				for {
					if len(outDP.Queue) == 0 {
						break
					}

					x := outDP.Remove()

					if item, okItem := x.(*Item); okItem {
						for _, i := range item.Splay() {
							items = append(items, i)
						}
					}

					if promise, okPromise := x.(*TaskPromise); okPromise {
						for _, p := range promise.Splay() {
							newTask := NewTaskFromPromise(p, w)
							tasks = append(tasks, newTask)
						}
					}
				}
			}
		}
	}

	return items, nil
}
