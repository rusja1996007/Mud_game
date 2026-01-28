package todo

import "errors"

var ErrTaskNotFound = errors.New("task not found")           //task - задача
var ErrTaskAlreadyExists = errors.New("task already exists") //попытка создать существующую задачу
