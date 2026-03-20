package psql

import (
	"time"
)

func (db *DataBase) Insert( 
	title string, 
	description string, 
	isCompleted bool, 
	createdAt time.Time, 
	completedAt *time.Time) error {

	query := `INSERT INTO Tasks (
	title, description, isCompleted, createdAt, completedAt
	)
	VALUES($1, $2, $3, $4, $5)
	RETURNING title, description, isCompleted, createdAt, completedAt
	`
	_, err := db.Conn.Exec(db.Ctx, query, title, description, isCompleted, createdAt, completedAt)
	return err
}