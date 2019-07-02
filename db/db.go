package db

import (
	"database/sql"
)

//RegisterUser adds user to DB
func RegisterUser(DB *sql.DB, login, userFIO string) error {
	stmt := "INSERT INTO profiles (id, username, user_fio) VALUES (DEFAULT, $1, $2) RETURNING id"
	var userID int64
	err := DB.QueryRow(stmt, login, userFIO).Scan(&userID)
	if err != nil {
		return err
	}

	return nil
}

func prepareStmt(db *sql.DB, stmt string) (*sql.Stmt, error) {
	res, err := db.Prepare(stmt)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func InitStmts() {
}
