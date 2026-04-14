package psql

import (
)

func (db *DataBase) Delete(id int, userId int) error {
	query := `DELETE FROM Tasks
	WHERE id = $1 AND user_id = $2
	`
	_, err := db.Conn.Exec(db.Ctx, query, id, userId)
	return err
}
