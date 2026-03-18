package psql

import (
	"context"

	"github.com/jackc/pgx/v5"
)

func Select(ctx context.Context, conn *pgx.Conn) error {
	query := `SELECT title, description, isCompleted, createdAt, completedAt
	FROM Tasks
	`
	_, err := conn.Exec(ctx, query)
	return err
}