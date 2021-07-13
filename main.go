package main

import (
	"fmt"
	"net/http"

	"github.com/davecgh/go-spew/spew"
)

type DataPipe struct {
	Done  bool
	Queue []interface{}
}

type AbstractAction struct {
	Inputs     map[string]DataPipe
	Outputs    map[string]DataPipe
	CanFail    bool
	ExpectMany bool
}

type Action interface {
	Run() bool
}

type HTTPAction struct {
	AbstractAction
}

type Task struct {
	Inputs  map[string]DataPipe
	Outputs map[string]DataPipe
}

type Workflow struct {
	Name    string
	Version string
	Tasks   []Task
}

func main() {
	fmt.Println("SpiderSwarm")
	httpClient := &http.Client{}
	spew.Dump(httpClient)

	w := Workflow{}
	spew.Dump(w)
}
