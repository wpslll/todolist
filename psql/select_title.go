package psql

func (db *DataBase) Select_title(title string, userId int) (TaskDto, error) {
	query := `SELECT title, description, isCompleted, createdAt, completedAt
	FROM Tasks
	WHERE title = $1 AND user_id = $2
	`
	var task TaskDto
	err := db.Conn.QueryRow(db.Ctx, query, title, userId).Scan(&task.Title, &task.Description, &task.IsCompleted, &task.CreatedAt, &task.CompletedAt)
	return task, err
}