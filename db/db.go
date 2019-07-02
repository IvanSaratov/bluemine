package db

import (
	"database/sql"

	"github.com/IvanSaratov/bluemine/data"
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

//GetUserInfo gets user info from DB
func GetUserInfo(DB *sql.DB, login string) (data.User, error) {
	var user data.User

	stmt := "SELECT * FROM profiles WHERE username = $1"
	err := DB.QueryRow(stmt, login).Scan(&user.UserID, &user.UserName, &user.UserFIO, &user.UserisAdmin, &user.UserRate)
	if err != nil {
		return user, err
	}

	return user, nil
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
