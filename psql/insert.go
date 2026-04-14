package psql

import (
	"time"
)

func (db *DataBase) Insert( 
	title string, 
	description string, 
	isCompleted bool, 
	createdAt time.Time, 
	completedAt *time.Time,
	userId int) (TaskDto, error) {

	query := `INSERT INTO Tasks (
	title, description, isCompleted, createdAt, completedAt
	)
	VALUES($1, $2, $3, $4, $5, $6)
	RETURNING id, title, description, isCompleted, createdAt, completedAt
	`
	var task TaskDto
	err := db.Conn.QueryRow(db.Ctx, query, title, description, isCompleted, createdAt, completedAt, userId).Scan(&task.Id, &task.Title, &task.Description, &task.IsCompleted, &task.CreatedAt, &task.CompletedAt)
	return task, err
}