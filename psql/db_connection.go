package psql

import (
	"context"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

type DataBase struct {
	Conn *pgx.Conn
	Ctx context.Context
}

func NewDB(conn *pgx.Conn, ctx context.Context) DataBase{
	return DataBase {
		Conn: conn,
		Ctx: ctx,
	}
}

func CheckConnection(ctx context.Context) (*pgx.Conn, error){
	if err := godotenv.Load(); err != nil {
		panic(err)
	}
	db_url := os.Getenv("DB_URL_DOCKER")
	return pgx.Connect(ctx, db_url)
}
