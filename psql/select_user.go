package psql

func (db *DataBase) SelectUser(login, password string) (int, error) {
	query := `SELECT id
	FROM Users
	WHERE login = $1 AND password = $2
	`
	var id int
	err := db.Conn.QueryRow(db.Ctx, query, login, password).Scan(&id)
	if err != nil {
		return -1, err
	}
	return id, nil
}

func (db *DataBase) SelectUserId(id int) error {
	query := `SELECT login
	FROM Users
	WHERE id = $1
	`
	var login string
	return db.Conn.QueryRow(db.Ctx, query, id).Scan(&login)
}