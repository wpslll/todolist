package psql

import (
	"context"

	"github.com/jackc/pgx/v5"
)

func Select_title(ctx context.Context, conn *pgx.Conn, title string) error {
	query := `SELECT title, description, isCompleted, createdAt, completedAt
	FROM Tasks
	WHERE ID = $1
	`
	_, err := conn.Exec(ctx, query, title)
	return err
}