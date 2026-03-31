package http

import (
	"backend/middleware"
	"errors"
	"net/http"

	"github.com/gorilla/mux"
)

type HttpServer struct {
	httpHandlers *HttpHandlers
}

func NewHttpServer(httpHandler *HttpHandlers) *HttpServer {
	return &HttpServer{
		httpHandlers: httpHandler,
	}
}

func (s *HttpServer) StartServer() error {
	router := mux.NewRouter()
	router.Path("/tasks").Methods("POST").HandlerFunc(s.httpHandlers.HandleCreateTask)
	router.Path("/tasks").Methods("GET").Queries("completed", "false").HandlerFunc(middleware.AuthMiddleware(s.httpHandlers.HandleGetAllUncompletedTasks))
	router.Path("/tasks/{title}").Methods("GET").HandlerFunc(middleware.AuthMiddleware(s.httpHandlers.HandleGetTask))
	router.Path("/tasks").Methods("GET").HandlerFunc(middleware.AuthMiddleware(s.httpHandlers.HandleGetAllTasks))
	router.Path("/tasks/{id}").Methods("PATCH").HandlerFunc(middleware.AuthMiddleware(s.httpHandlers.HandleCompleteTask))
	router.Path("/tasks/{id}").Methods("DELETE").HandlerFunc(middleware.AuthMiddleware(s.httpHandlers.HandleDeleteTask))
	router.Path("/health").Methods("GET").HandlerFunc(s.httpHandlers.HandleHealth)
	router.Path("/register").Methods("POST").HandlerFunc(s.httpHandlers.HandleRegistration)
	router.Path("/auth").Methods("POST").HandlerFunc(s.httpHandlers.HandleAuth)
	router.Path("/refresh").Methods("POST").HandlerFunc(s.httpHandlers.HandleRefresh)
	if err := http.ListenAndServe(":9091", router); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err
	}
	return nil
}
