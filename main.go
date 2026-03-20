package main

import (
	"backend/http"
	"backend/psql"
	"backend/todo"
	"context"
	"fmt"
)

func main() {
	ctx := context.Background()
	conn, err := psql.CheckConnection(ctx)
	if err != nil {
		fmt.Println("Failed to connect to database: ", err)
	}
	db := psql.NewDB(conn, ctx)
	todoList := todo.NewList(db)
	httpHandler := http.NewHttpHandler(todoList)
	httpserver := http.NewHttpServer(httpHandler)
	if err := httpserver.StartServer(); err != nil {
		fmt.Println("Failed to start server: ", err)
	}
}