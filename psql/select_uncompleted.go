package psql

import "fmt"

func (db *DataBase) Select_uncompleted(userId int) (map[string]TaskDto, error) {
	query := `SELECT title, description, isCompleted, createdAt, completedAt
	FROM Tasks
	WHERE isCompleted IS NOT TRUE AND user_id = $1
	RETURNING title, description, isCompleted, createdAt, completedAt
	`
	list := make(map[string]TaskDto)
	rows, err := db.Conn.Query(db.Ctx, query, userId)
	if err != nil {
		fmt.Println("Failed to get rows from db: ", err)
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var task TaskDto
		if err := rows.Scan(&task.Title, &task.Description, &task.IsCompleted, &task.CreatedAt, &task.CompletedAt); err != nil {
			fmt.Println("Failed to scan into a variable: ", err)
			return nil, err
		}
		list[task.Title] = task
	}
	return list, nil
}
