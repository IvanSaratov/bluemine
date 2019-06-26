package db

import (
	"database/sql"
)

func prepareStmt(db *sql.DB, stmt string) (*sql.Stmt, error) {
	res, err := db.Prepare(stmt)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func InitStmts() {
}
