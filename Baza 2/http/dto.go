package http

import (
	"encoding/json"
	"errors"
	"time"
)

type TaskDTO struct {
	Title    string
	Opisanie string
}

func (t TaskDTO) ValidateForCreate() error {
	if t.Title == "" {
		return errors.New("title пустой")
	}
	if t.Opisanie == "" {
		return errors.New("opisanie пустое")
	}
	return nil

}

type ErrorDTO struct {
	Message string
	Time    time.Time
}

func (e ErrorDTO) ToString() string {
	b, err := json.MarshalIndent(e, "", "    ")
	if err != nil {
		panic(err)
	}
	return string(b)
}

/*
"complete": true
*/
type CompleteTaskDTO struct {
	Complete bool
}
