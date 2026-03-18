package psql

import (
	"context"

	"github.com/jackc/pgx/v5"
)

func Delete(ctx context.Context, conn *pgx.Conn, id int) error {
	query := `DELETE FROM Tasks
	WHERE id = $1
	`
	_, err := conn.Exec(ctx, query, id)
	return err
}