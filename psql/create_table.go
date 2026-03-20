package psql

import (
)

func (db *DataBase) CreateTable() error{
	query := `CREATE TABLE Tasks (
		id SERIAL PRIMARY KEY,
		title VARCHAR(200) NOT NULL,
		description VARCHAR(1000),
		isCompleted BOOLEAN NOT NULL,
		createdAt TIMESTAMP NOT NULL,
		completedAt TIMESTAMP
	);`
	_, err := db.Conn.Exec(db.Ctx, query)
	return err
}