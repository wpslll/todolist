package psql

import (
)

func (db *DataBase) Delete(title string) error {
	query := `DELETE FROM Tasks
	WHERE title = $1
	`
	_, err := db.Conn.Exec(db.Ctx, query, title)
	return err
}
