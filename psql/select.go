package psql


func (db *DataBase) Select(id int) (map[string]TaskDto, error) {
	query := `SELECT id, title, description, isCompleted, createdAt, completedAt
	FROM Tasks
	WHERE user_id = $1
	`
	list := make(map[string]TaskDto)
	rows, err := db.Conn.Query(db.Ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var task TaskDto
		if err := rows.Scan(&task.Id, &task.Title, &task.Description, &task.IsCompleted, &task.CreatedAt, &task.CompletedAt); err != nil {
			panic(err)
		}
		list[task.Title] = task
	}
	return list, err
}