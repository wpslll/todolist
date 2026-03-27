package psql

import (
	"time"
)

func (db *DataBase)Complete_Uncomplete(title string, isCompleted bool, completedAt *time.Time) error{
	query := `UPDATE Tasks
	SET isCompleted = $2, completedAt = $3
	WHERE title = $1
	RETURNING title, description, isCompleted, createdAt, completedAt
	`
	_, err := db.Conn.Exec(db.Ctx, query, title, isCompleted, completedAt)
	return err
}