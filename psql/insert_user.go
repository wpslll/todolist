package psql

func (db *DataBase) InsertUser(login string, password string) error {
	query := `INSERT INTO Users (
	login, password
	)
	VALUES($1, $2)
	RETURNING login, password
	`
	db.Conn.Exec(db.Ctx, query, login, password)
	return nil
}