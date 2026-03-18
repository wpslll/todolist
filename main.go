package main

import (
	"backend/http"
	"backend/psql"
	"backend/todo"
	"context"
	"fmt"
)

func main() {
	todoList := todo.NewList()
	httpHandler := http.NewHttpHandler(todoList)
	httpserver := http.NewHttpServer(httpHandler)
	ctx := context.Background()
	psql.CheckConnection(ctx)
	if err := httpserver.StartServer(); err != nil {
		fmt.Println("Failed to start server: ", err)
	}
}