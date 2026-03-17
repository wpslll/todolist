package main

import (
	"backend/http"
	"backend/todo"
	"fmt"
)

func main() {
	todoList := todo.NewList()
	httpHandler := http.NewHttpHandler(todoList)
	httpserver := http.NewHttpServer(httpHandler)
	if err := httpserver.StartServer(); err != nil {
		fmt.Println("Failed to start server: ", err)
	}
}