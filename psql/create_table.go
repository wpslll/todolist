package psql

import (
	"context"

	"github.com/jackc/pgx/v5"
)

func CreateTable(conn *pgx.Conn, ctx context.Context) error{
	query := `CREATE TABLE Tasks (
		id SERIAL PRIMARY KEY,
		title VARCHAR(200) NOT NULL,
		description VARCHAR(1000),
		isCompleted BOOLEAN NOT NULL,
		createdAt TIMESTAMP NOT NULL,
		completedAt TIMESTAMP
	);`
	_, err := conn.Exec(ctx, query)
	return err
}