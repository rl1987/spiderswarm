package spsw

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/google/uuid"
)

type TaskResult struct {
	UUID              string
	JobUUID           string
	TaskUUID          string
	ScheduledTaskUUID string
	Succeeded         bool
	Error             error
	OutputDataChunks  map[string][]*DataChunk
}

func NewTaskResult(jobUUID string, taskUUID string, scheduledTaskUUID string, succeeded bool, err error) *TaskResult {
	return &TaskResult{
		UUID:              uuid.New().String(),
		JobUUID:           jobUUID,
		TaskUUID:          taskUUID,
		ScheduledTaskUUID: scheduledTaskUUID,
		Succeeded:         succeeded,
		Error:             err,
		OutputDataChunks:  map[string][]*DataChunk{},
	}
}

func NewTaskResultFromJSON(raw []byte) *TaskResult {
	taskResult := &TaskResult{}

	buffer := bytes.NewBuffer(raw)
	decoder := json.NewDecoder(buffer)

	err := decoder.Decode(taskResult)
	if err != nil {
		return nil
	}

	return taskResult
}

func (tr *TaskResult) String() string {
	return fmt.Sprintf("<TaskResult %s JobUUID: %s, TaskUUID: %s, ScheduledTaskUUID: %s, Succeeded: %v, Error: %v, OutputDataChunks: %v>",
		tr.UUID, tr.JobUUID, tr.TaskUUID, tr.ScheduledTaskUUID, tr.Succeeded, tr.Error, tr.OutputDataChunks)
}

func (tr *TaskResult) EncodeToJSON() []byte {
	buffer := bytes.NewBuffer([]byte{})
	encoder := json.NewEncoder(buffer)

	encoder.Encode(tr)

	bytes, _ := ioutil.ReadAll(buffer)

	return bytes
}

func (tr *TaskResult) addOutputDataChunk(outputName string, chunk *DataChunk) {
	if tr.OutputDataChunks[outputName] == nil {
		tr.OutputDataChunks[outputName] = []*DataChunk{chunk}
	} else {
		tr.OutputDataChunks[outputName] = append(tr.OutputDataChunks[outputName], chunk)
	}
}

func (tr *TaskResult) AddOutputItem(outputName string, item *Item) {
	chunk, _ := NewDataChunk(item)

	tr.addOutputDataChunk(outputName, chunk)
}

func (tr *TaskResult) AddOutputTaskPromise(outputName string, promise *TaskPromise) {
	chunk, _ := NewDataChunk(promise)

	tr.addOutputDataChunk(outputName, chunk)
}
