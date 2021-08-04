package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewActionFromTemplate(t *testing.T) {
	workflow := &Workflow{}
	jobUUID := "B3FB46CC-E3AC-41FB-AC88-A15B23C2F299"

	actionTempl1 := &ActionTemplate{StructName: "HTTPAction"}
	action1, ok1 := NewActionFromTemplate(actionTempl1, workflow, jobUUID).(*HTTPAction)

	assert.True(t, ok1)
	assert.NotNil(t, action1)

	actionTempl2 := &ActionTemplate{
		StructName: "XPathAction",
		ConstructorParams: map[string]interface{}{
			"xpath":      "//title",
			"expectMany": false,
		},
	}
	action2, ok2 := NewActionFromTemplate(actionTempl2, workflow, jobUUID).(*XPathAction)

	assert.True(t, ok2)
	assert.NotNil(t, action2)

	actionTempl3 := &ActionTemplate{StructName: "FieldJoinAction"}
	action3, ok3 := NewActionFromTemplate(actionTempl3, workflow, jobUUID).(*FieldJoinAction)

	assert.True(t, ok3)
	assert.NotNil(t, action3)

	actionTempl4 := &ActionTemplate{StructName: "TaskPromiseAction"}
	action4, ok4 := NewActionFromTemplate(actionTempl4, workflow, jobUUID).(*TaskPromiseAction)

	assert.True(t, ok4)
	assert.NotNil(t, action4)

	actionTempl5 := &ActionTemplate{StructName: "UTF8DecodeAction"}
	action5, ok5 := NewActionFromTemplate(actionTempl5, workflow, jobUUID).(*UTF8DecodeAction)

	assert.True(t, ok5)
	assert.NotNil(t, action5)

	actionTempl6 := &ActionTemplate{StructName: "UTF8EncodeAction"}
	action6, ok6 := NewActionFromTemplate(actionTempl6, workflow, jobUUID).(*UTF8EncodeAction)

	assert.True(t, ok6)
	assert.NotNil(t, action6)
}
