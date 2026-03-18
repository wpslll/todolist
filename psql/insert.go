package psql

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
)

func Insert(
	ctx context.Context, 
	conn *pgx.Conn, 
	title string, 
	description string, 
	isCompleted bool, 
	createdAt time.Time, 
	completedAt time.Time) error {
	query := `INSERT INTO Tasks (
	title, description, isCompleted, createdAt, completedAt
	)
	VALUES($1, $2, $3, $4, $5)
	`
	_, err := conn.Exec(ctx, query, title, description, isCompleted, createdAt, completedAt)
	return err
}