package status

import (
	"fmt"
	"sync"
)

var lock = &sync.Mutex{}

type single struct {
	isRunning bool
	messages  []message
}

type message struct {
	msgType    string
	msgSubject string
	msg        string
}

var singleInstance *single

func getInstance() *single {
	if singleInstance == nil {
		lock.Lock()
		defer lock.Unlock()
		if singleInstance == nil {
			singleInstance = &single{isRunning: false, messages: make([]message, 0)}
		}
	}
	return singleInstance
}

type Status struct {
}

func GetStatus() Status {
	return Status{}
}

func (s Status) StartProcess() error {
	single := getInstance()
	lock.Lock()
	defer lock.Unlock()
	if single.isRunning {
		return fmt.Errorf("already running")
	} else {
		single.isRunning = true
	}
	return nil
}
func (s Status) EndProcess() error {
	single := getInstance()
	lock.Lock()
	defer lock.Unlock()
	if !single.isRunning {
		return fmt.Errorf("not running")
	} else {
		single.isRunning = false
	}
	return nil
}
func (s Status) Success(book string, msg string) {
	single := getInstance()
	lock.Lock()
	defer lock.Unlock()
	single.messages = append(single.messages, message{msgType: "success", msgSubject: book, msg: msg})

}
func (s Status) Error(book string, msg string) {
	single := getInstance()
	lock.Lock()
	defer lock.Unlock()
	single.messages = append(single.messages, message{msgType: "error", msgSubject: book, msg: msg})

}
