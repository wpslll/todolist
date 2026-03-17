package todo

import "errors"

var ServerConnectionError = errors.New("Failed to connect http server: ")
var TaskNotFound = errors.New("Failed to find task: ")
var TaskAlreadyExists = errors.New("Task already exists: ")