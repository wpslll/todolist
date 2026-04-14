package psql

func (db *DataBase) UpdateTask(id int, userId int, title string, description string) (TaskDto, error) {
	query := `UPDATE Tasks
	SET title = $3, description = $4
	WHERE id = $1 AND user_id = $2
	RETURNING id, title, description, isCompleted, createdAt, completedAt
	`
	var task TaskDto
	err := db.Conn.QueryRow(db.Ctx, query, id, userId, title, description).Scan(&task.Id, &task.Title, &task.Description, &task.IsCompleted, &task.CreatedAt, &task.CompletedAt)
	return task, err
}