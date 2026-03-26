package main

import (
	"backend/http"
	"backend/psql"
	"backend/todo"
	"context"
	"fmt"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
)

func main() {
	ctx := context.Background()
	conn, err := psql.CheckConnection(ctx)
	if err != nil {
		fmt.Println("Failed to connect to database: ", err)
	}
	if err := godotenv.Load(); err != nil {
		fmt.Println("Failed to load env: ", err)
	}
	db_url := os.Getenv("DB_URL_DOCKER")
	db := psql.NewDB(conn, ctx)
	m, err := migrate.New("file://migrations", db_url)
	if err != nil {
		fmt.Println("Failed to create a migration: ", err)
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		fmt.Println("Failed to migrate: ", err)
	}
	todoList := todo.NewList(db)
	httpHandler := http.NewHttpHandler(todoList)
	httpserver := http.NewHttpServer(httpHandler)
	if err := httpserver.StartServer(); err != nil {
		fmt.Println("Failed to start server: ", err)
	}
}
