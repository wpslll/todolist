package psql

import (
)

func (db *DataBase) Delete(id int) error {
	query := `DELETE FROM Tasks
	WHERE id = $1
	`
	_, err := db.Conn.Exec(db.Ctx, query, id)
	return err
}
