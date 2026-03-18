package psql

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
)

func Complete_Uncomplete(ctx context.Context, conn *pgx.Conn, id int, isCompleted bool, completedAt time.Time) error{
	query := `UPDATE Tasks
	SET isCompleted = $2, completedAt = $3
	WHERE id = $1
	`
	_, err := conn.Exec(ctx, query, id, isCompleted, completedAt)
	return err
}