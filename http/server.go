package http

import (
	"backend/middleware"
	"errors"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
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
	router.Path("/tasks").Methods("GET").Queries("completed", "false").HandlerFunc(middleware.AuthMiddleware(s.httpHandlers.HandleGetAllUncompletedTasks, *s.httpHandlers.logger))
	router.Path("/tasks/{title}").Methods("GET").HandlerFunc(middleware.AuthMiddleware(s.httpHandlers.HandleGetTask, *s.httpHandlers.logger))
	router.Path("/tasks").Methods("GET").HandlerFunc(middleware.AuthMiddleware(s.httpHandlers.HandleGetAllTasks, *s.httpHandlers.logger))
	router.Path("/tasks/{id}").Methods("PATCH").HandlerFunc(middleware.AuthMiddleware(s.httpHandlers.HandleCompleteTask, *s.httpHandlers.logger))
	router.Path("/tasks/{id}").Methods("DELETE").HandlerFunc(middleware.AuthMiddleware(s.httpHandlers.HandleDeleteTask, *s.httpHandlers.logger))
	router.Path("/tasks/{id}").Methods("UPDATE").HandlerFunc(middleware.AuthMiddleware(s.httpHandlers.HandleUpdateTask, *s.httpHandlers.logger))
	router.Path("/health").Methods("GET").HandlerFunc(s.httpHandlers.HandleHealth)
	router.Path("/register").Methods("POST").HandlerFunc(s.httpHandlers.HandleRegistration)
	router.Path("/auth").Methods("POST").HandlerFunc(s.httpHandlers.HandleAuth)
	router.Path("/refresh").Methods("POST").HandlerFunc(s.httpHandlers.HandleRefresh)
	router.Methods("OPTIONS").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
})
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:3000"},
		AllowedMethods: []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
		Debug: false,
	})
	handler := c.Handler(router)
	if err := http.ListenAndServe(":9091", handler); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err
	}
	return nil
}
