// Copyright © 2019 National Library of Norway
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package importutil

import (
	"fmt"
	"os"
	"time"
)

type executor struct {
	Count      int
	Success    int
	Failed     int
	Processor  func(value interface{}, state *State) error
	ErrHandler func(state *StateVal)
	LogHandler func(state *StateVal)

	threadCount  int
	dataChan     chan *StateVal
	completeChan chan *processorResponse
	errorChan    chan *StateVal
	errsHandled  int
}

type StateVal struct {
	*State
	Val interface{}
}

type processorResponse struct {
	count   int
	success int
	failed  int
}

func NewExecutor(threadCount int, proc func(value interface{}, state *State) error, errorHandler func(state *StateVal)) *executor {
	e := &executor{
		threadCount:  threadCount,
		Processor:    proc,
		ErrHandler:   errorHandler,
		dataChan:     make(chan *StateVal),
		completeChan: make(chan *processorResponse),
		errorChan:    make(chan *StateVal),
	}

	go e.handleError()

	for i := 0; i < e.threadCount; i++ {
		go e.execute()
	}

	return e
}

func (e *executor) Do(state *State, val interface{}) {
	if val != nil {
		e.dataChan <- &StateVal{state, val}
		e.Count++
		e.printProgress()
	}
}

func (e *executor) Finish() {
	close(e.dataChan)
	subCompleted := 0
	for i := range e.completeChan {
		e.Success += i.success
		e.Failed += i.failed
		subCompleted++
		if subCompleted >= e.threadCount {
			break
		}
	}
	for e.errsHandled < e.Failed {
		time.Sleep(time.Second)
	}
}

func (e *executor) execute() {
	res := &processorResponse{}

	for request := range e.dataChan {
		res.count++
		err := e.Processor(request.Val, request.State)
		if err != nil {
			request.State.err = err
			e.errorChan <- request
			res.failed++
		} else {
			res.success++
		}
	}
	e.completeChan <- res
}

func (e *executor) handleError() {
	for err := range e.errorChan {
		e.ErrHandler(err)
		e.errsHandled++
	}
}

func (e *executor) printProgress() {
	// Print progress
	_, _ = fmt.Fprint(os.Stderr, ".")
	if e.Count%100 == 0 {
		_, _ = fmt.Fprintln(os.Stderr, e.Count)
	}
}
