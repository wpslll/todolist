package http

import (
	"encoding/json"
	"errors"
	"time"
)
type CompleteTaskDTO struct {
	Complete bool
}

type TaskDTO struct {
	Title       string
	Description string
}

type UserDTO struct {
	Login 	 string
	Password string
}

func (t TaskDTO) ValidateForCreate() error {
	if t.Title == "" {
		return errors.New("Title is empty")
	}
	if t.Description == "" {
		return errors.New("Description is empty")
	}
	return nil
}

type ErrorDTO struct {
	Message string
	Time    time.Time
}

func (e ErrorDTO) ToString() string {
	b, err := json.MarshalIndent(e, "", "	")
	if err != nil {
		panic(err)
	}
	return string(b)
}

func (u *UserDTO) ValidateUser() error {
	if u.Login == "" {
		return errors.New("Username is empty")
	}
	if u.Password == "" {
		return errors.New("Password is empty")
	}
	return nil
}