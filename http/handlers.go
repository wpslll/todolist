package http

import (
	"backend/middleware"
	"backend/todo"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type HttpHandlers struct {
	todoList *todo.List
	logger *zap.Logger
}

func NewHttpHandler(todolist *todo.List, log *zap.Logger) *HttpHandlers {
	return &HttpHandlers{
		todoList: todolist,
		logger: log,
	}
}

func (h *HttpHandlers) HandleCreateTask(w http.ResponseWriter, r *http.Request) {
	var taskDTO TaskDTO
	if err := json.NewDecoder(r.Body).Decode(&taskDTO); err != nil {
		errDTO := ErrorDTO{
			Message: err.Error(),
			Time:    time.Now(),
		}
		h.logger.Error("Failed to Unmarshal json")
		http.Error(w, errDTO.ToString(), http.StatusBadRequest)
		return
	}
	if err := taskDTO.ValidateForCreate(); err != nil {
		errDTO := ErrorDTO{
			Message: err.Error(),
			Time:    time.Now(),
		}
		h.logger.Error("Failed to validate task")
		http.Error(w, errDTO.ToString(), http.StatusBadRequest)
		return
	}
	h.logger.Info("Creating task")
	todoTask := todo.NewTask(taskDTO.Title, taskDTO.Description)
	if err := h.todoList.AddTask(todoTask); err != nil {
		errDTO := ErrorDTO{
			Message: err.Error(),
			Time:    time.Now(),
		}
		h.logger.Error("Failed to add task")
		if errors.Is(err, todo.TaskAlreadyExists) {
			http.Error(w, errDTO.ToString(), http.StatusConflict)
		} else {
			http.Error(w, errDTO.ToString(), http.StatusInternalServerError)
		}
		return
	}
	b, err := json.MarshalIndent(todoTask, "", "    ")
	if err != nil {
		h.logger.Error("Failed to marshal task")
		panic(err)
	}
	w.WriteHeader(http.StatusCreated)
	if _, err := w.Write(b); err != nil {
		h.logger.Error("Failed to write response body")
		fmt.Println("Failed to write http response: ", err)
		return
	}
}

func (h *HttpHandlers) HandleGetTask(w http.ResponseWriter, r *http.Request) {
	title := mux.Vars(r)["title"]
	h.logger.Info("Getting task")
	task, err := h.todoList.GetTask(title)
	if err != nil {
		errDTO := ErrorDTO{
			Message: err.Error(),
			Time:    time.Now(),
		}
		h.logger.Error("Failed to get task")
		if errors.Is(err, todo.TaskNotFound) {
			http.Error(w, errDTO.ToString(), http.StatusNotFound)
		} else {
			http.Error(w, errDTO.ToString(), http.StatusInternalServerError)
		}
		return
	}
	b, err := json.MarshalIndent(task, "", "    ")
	if err != nil {
		h.logger.Error("Failed to marshal json")
		panic(err)
	}
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(b); err != nil {
		h.logger.Error("Failed to write response body")
		fmt.Println("Failed to write response: ", err)
		return
	}
}

func (h *HttpHandlers) HandleGetAllTasks(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.todoList.ListTasks()
	if err != nil {
		h.logger.Error("Failed to get all tasks")
		fmt.Println("Failed to get tasks from db: ", err)
	}
	b, err := json.MarshalIndent(tasks, "", "    ")
	if err != nil {
		h.logger.Error("Failed to marshal json")
		panic(err)
	}
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(b); err != nil {
		h.logger.Error("Failed to write response body")
		fmt.Println("Failed to write response: ", err)
		return
	}
}

func (h *HttpHandlers) HandleGetAllUncompletedTasks(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("Getting all uncompleted tasks")
	uncompletedTasks, err := h.todoList.ListUncompletedTasks()
	if err != nil {
		h.logger.Error("Failed to get all uncompleted tasks")
		fmt.Println("Failed to get uncompleted tasks from db: ", err)
		return
	}
	b, err := json.MarshalIndent(uncompletedTasks, "", "    ")
	if err != nil {
		h.logger.Error("Failed to marshal json")
		panic(err)
	}
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(b); err != nil {
		h.logger.Error("Failed to write response body")
		fmt.Println("Failed to write response: ", err)
		return
	}
}

func (h *HttpHandlers) HandleCompleteTask(w http.ResponseWriter, r *http.Request) {
	var completeDTO CompleteTaskDTO
	h.logger.Info("Reading from req body")
	if err := json.NewDecoder(r.Body).Decode(&completeDTO); err != nil {
		errDTO := ErrorDTO{
			Message: err.Error(),
			Time:    time.Now(),
		}
		h.logger.Error("Failed to decode req body")
		http.Error(w, errDTO.ToString(), http.StatusBadRequest)
		return
	}

	title := mux.Vars(r)["title"]

	var (
		changedTask todo.Task
		err         error
	)
	h.logger.Info("Completing or Uncompleting task")
	if completeDTO.Complete {
		changedTask, err = h.todoList.CompleteTask(title)
	} else {
		changedTask, err = h.todoList.UncompleteTask(title)
	}

	if err != nil {
		errDTO := ErrorDTO{
			Message: err.Error(),
			Time:    time.Now(),
		}
		h.logger.Error("Failed to complete or uncomplete task")
		if errors.Is(err, todo.TaskNotFound) {
			http.Error(w, errDTO.ToString(), http.StatusNotFound)
		} else {
			http.Error(w, errDTO.ToString(), http.StatusInternalServerError)
		}
		return
	}
	b, err := json.MarshalIndent(changedTask, "", "    ")

	if _, err := w.Write(b); err != nil {
		h.logger.Error("Failed to write response body")
		fmt.Println("Failed to write response: ", err)
		return
	}
}

func (h *HttpHandlers) HandleDeleteTask(w http.ResponseWriter, r *http.Request) {
	title := mux.Vars(r)["title"]
	h.logger.Info("Deleting task")
	if err := h.todoList.DeleteTask(title); err != nil {
		errDTO := ErrorDTO{
			Message: err.Error(),
			Time:    time.Now(),
		}
		h.logger.Error("Failed to delete task")
		if errors.Is(err, todo.TaskNotFound) {
			http.Error(w, errDTO.ToString(), http.StatusNotFound)
		} else {
			http.Error(w, errDTO.ToString(), http.StatusInternalServerError)
		}
		return
	}
}

func (h *HttpHandlers) HandleHealth(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("Checking server health")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func (h *HttpHandlers) HandleRegistration(w http.ResponseWriter, r *http.Request) {
	var userdto UserDTO
	h.logger.Info("Reading from req body")
	if err := json.NewDecoder(r.Body).Decode(&userdto); err != nil {
		errdto := ErrorDTO{
			Message: err.Error(),
			Time:    time.Now(),
		}
		h.logger.Error("Failed to decode req body")
		http.Error(w, errdto.ToString(), http.StatusBadRequest)
		return
	}
	if err := userdto.ValidateUser(); err != nil {
		errdto := ErrorDTO{
			Message: err.Error(),
			Time:    time.Now(),
		}
		h.logger.Error("Failed to validate user")
		http.Error(w, errdto.ToString(), http.StatusBadRequest)
		return
	}
	h.logger.Info("Creating new user")
	if err := h.todoList.NewUser(userdto.Login, userdto.Password); err != nil {
		errdto := ErrorDTO{
			Message: err.Error(),
			Time:    time.Now(),
		}
		h.logger.Error("Failed to create new user")
		http.Error(w, errdto.ToString(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (h *HttpHandlers) HandleAuth(w http.ResponseWriter, r *http.Request) {
	var userdto UserDTO
	h.logger.Info("Reading from req body")
	if err := json.NewDecoder(r.Body).Decode(&userdto); err != nil {
		errdto := ErrorDTO{
			Message: err.Error(),
			Time:    time.Now(),
		}
		h.logger.Error("Failed to decode req body")
		http.Error(w, errdto.ToString(), http.StatusBadRequest)
		return
	}
	if err := userdto.ValidateUser(); err != nil {
		errdto := ErrorDTO{
			Message: err.Error(),
			Time:    time.Now(),
		}
		h.logger.Error("Failed to validate user")
		http.Error(w, errdto.ToString(), http.StatusBadRequest)
		return
	}
	h.logger.Info("Finding user")
	id, err := h.todoList.FindUser(userdto.Login, userdto.Password)
	if err != nil {
		errdto := ErrorDTO{
			Message: err.Error(),
			Time:    time.Now(),
		}
		h.logger.Error("Failed to find user")
		http.Error(w, errdto.ToString(), http.StatusBadRequest)
		return
	}
	h.logger.Info("Creating tokens for user")
	accessToken, refreshToken, err := middleware.CreateToken(id)
	if err != nil {
		errdto := ErrorDTO{
			Message: err.Error(),
			Time:    time.Now(),
		}
		h.logger.Error("Failed to create tokens")
		http.Error(w, errdto.ToString(), http.StatusInternalServerError)
		return
	}
	h.logger.Info("Setting cookies")
	http.SetCookie(w, &http.Cookie{
		Name:     "accessToken",
		Value:    accessToken,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
		MaxAge:   15 * 60,
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "refreshToken",
		Value:    refreshToken,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteStrictMode,
		Path:     "/refresh",
		MaxAge:   24 * 60 * 60,
	})
}

func (h *HttpHandlers) HandleRefresh(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("Reading from cookie")
	cookie, err := r.Cookie("refreshToken")
	if err != nil {
		errdto := ErrorDTO{
			Message: err.Error(),
			Time:    time.Now(),
		}
		h.logger.Error("Failed to read from cookie")
		http.Error(w, errdto.ToString(), http.StatusUnauthorized)
		return
	}
	tokenString := cookie.Value
	h.logger.Info("Refreshing tokens")
	accessToken, refreshToken, id, err := middleware.Refresh(tokenString)
	if err != nil {
		errdto := ErrorDTO{
			Message: err.Error(),
			Time:    time.Now(),
		}
		h.logger.Error("Failed to refresh tokens")
		http.Error(w, errdto.ToString(), http.StatusInternalServerError)
		return
	}
	h.logger.Info("Finding user")
	if err := h.todoList.FindUserId(id); err != nil {
		errdto := ErrorDTO{
			Message: err.Error(),
			Time:    time.Now(),
		}
		h.logger.Error("Failed to find user")
		http.Error(w, errdto.ToString(), http.StatusInternalServerError)
		return
	}
	h.logger.Info("Setting cookies")
	http.SetCookie(w, &http.Cookie{
		Name:     "accessToken",
		Value:    accessToken,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
		MaxAge:   15 * 60,
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "refreshToken",
		Value:    refreshToken,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteStrictMode,
		Path:     "/refresh",
		MaxAge:   24 * 60 * 60,
	})
}
