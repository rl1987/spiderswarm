package spsw

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewTaskResult(t *testing.T) {
	jobUUID := "A0"
	taskUUID := "B1"
	scheduledTaskUUID := "C2"
	succeeded := false
	err := errors.New("Test error")

	taskResult := NewTaskResult(jobUUID, taskUUID, scheduledTaskUUID, succeeded, err)

	assert.NotNil(t, taskResult)
	assert.Equal(t, jobUUID, taskResult.JobUUID)
	assert.Equal(t, taskUUID, taskResult.TaskUUID)
	assert.Equal(t, scheduledTaskUUID, taskResult.ScheduledTaskUUID)
	assert.Equal(t, succeeded, taskResult.Succeeded)
	assert.Equal(t, err, taskResult.Error)

	assert.NotNil(t, taskResult.OutputDataChunks)
	assert.Equal(t, 0, len(taskResult.OutputDataChunks))
}
