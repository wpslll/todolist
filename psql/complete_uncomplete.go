package psql

import (
	"time"
)

func (db *DataBase)Complete_Uncomplete(id int, isCompleted bool, completedAt *time.Time) (TaskDto, error){
	query := `UPDATE Tasks
	SET isCompleted = $2, completedAt = $3
	WHERE id = $1
	RETURNING id, title, description, isCompleted, createdAt, completedAt
	`
	var task TaskDto
	err := db.Conn.QueryRow(db.Ctx, query, id, isCompleted, completedAt).Scan(&task.Id, &task.Title, &task.Description, &task.IsCompleted, &task.CreatedAt, &task.CompletedAt)
	return task, err
}