package psql

import (
	"context"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

func CheckConnection(ctx context.Context) (*pgx.Conn, error){
	if err := godotenv.Load(); err != nil {
		panic(err)
	}
	db_url := os.Getenv("DB_URL_DOCKER")
	return pgx.Connect(ctx, db_url)
}
