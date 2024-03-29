package spsw

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
)

const StringMapUpdateActionInputOld = "StringMapUpdateActionInputOld"
const StringMapUpdateActionInputNew = "StringMapUpdateActionInputNew"
const StringMapUpdateActionInputOverridenValue = "StringMapUpdateActionInputOverridenValue"
const StringMapUpdateActionOutputUpdated = "StringMapUpdateActionOutputUpdated"
const StringMapUpdateActionOutputItem = "StringMapUpdateActionOutputItem"

type StringMapUpdateAction struct {
	AbstractAction
	OverrideKey string
	ItemName    string
}

func NewStringMapUpdateAction(overrideKey string) *StringMapUpdateAction {
	return &StringMapUpdateAction{
		AbstractAction: AbstractAction{
			CanFail:    false,
			ExpectMany: false,
			AllowedInputNames: []string{
				StringMapUpdateActionInputOld,
				StringMapUpdateActionInputNew,
				StringMapUpdateActionInputOverridenValue,
			},
			AllowedOutputNames: []string{
				StringMapUpdateActionOutputUpdated,
			},
			Inputs:  map[string]*DataPipe{},
			Outputs: map[string][]*DataPipe{},
			UUID:    uuid.New().String(),
		},
		OverrideKey: overrideKey,
	}
}

func NewStringMapUpdateActionFromTemplate(actionTempl *ActionTemplate) Action {
	action := NewStringMapUpdateAction("")

	action.Name = actionTempl.Name

	if _, ok := actionTempl.ConstructorParams["overrideKey"]; ok {
		action.OverrideKey = actionTempl.ConstructorParams["overrideKey"].StringValue
	}

	if _, ok := actionTempl.ConstructorParams["itemName"]; ok {
		action.ItemName = actionTempl.ConstructorParams["itemName"].StringValue
	}

	return action
}

func (smua *StringMapUpdateAction) String() string {
	return fmt.Sprintf("<StringMapUpdateAction %s Name: %s OverrideKey: %s>", smua.UUID, smua.Name, smua.OverrideKey)
}

func (smua *StringMapUpdateAction) Run() error {
	if smua.OverrideKey == "" {
		if smua.Inputs[StringMapUpdateActionInputOld] == nil || smua.Inputs[StringMapUpdateActionInputNew] == nil {
			return errors.New("Both inputs must be connected")
		}
	}

	if (smua.Outputs[StringMapUpdateActionOutputUpdated] == nil || len(smua.Outputs[StringMapUpdateActionOutputUpdated]) == 0) && (smua.Outputs[StringMapUpdateActionOutputItem] == nil || len(smua.Outputs[StringMapUpdateActionOutputItem]) == 0) {
		return errors.New("Output not connected")
	}

	updatedMap := map[string]string{}

	oldMap, ok1 := smua.Inputs[StringMapUpdateActionInputOld].Remove().(map[string]string)
	if !ok1 {
		return errors.New("Failed to get old cookies")
	}

	for key, value := range oldMap {
		updatedMap[key] = value
	}

	if smua.Inputs[StringMapUpdateActionInputNew] != nil {
		newMap, ok2 := smua.Inputs[StringMapUpdateActionInputNew].Remove().(map[string]string)
		if !ok2 {
			return errors.New("Failed to get new cookies")
		}

		for key, value := range newMap {
			updatedMap[key] = value
		}
	}

	if smua.OverrideKey != "" && smua.Inputs[StringMapUpdateActionInputOverridenValue] != nil {
		if valueStr, ok := smua.Inputs[StringMapUpdateActionInputOverridenValue].Remove().(string); ok {
			updatedMap[smua.OverrideKey] = valueStr
		}
	}

	if smua.Outputs[StringMapUpdateActionOutputUpdated] != nil {
		for _, output := range smua.Outputs[StringMapUpdateActionOutputUpdated] {
			output.Add(updatedMap)
		}
	}

	if smua.Outputs[StringMapUpdateActionOutputItem] != nil {
		for _, output := range smua.Outputs[StringMapUpdateActionOutputItem] {
			item := NewItem(smua.ItemName, "", "", "")

			for key, value := range updatedMap {
				item.SetField(key, value)
			}

			output.Add(item)
		}
	}

	return nil
}
