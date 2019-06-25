package db

import (
	"database/sql"
	"log"
)

var (
	Db *sql.DB
)

func prepareStmt(db *sql.DB, stmt string) *sql.Stmt{
	result, err := db.Prepare(stmt)
	if err != nil {
		log.Fatal("Could not prepare `" + stmt + "`: " + err.Error())
	}

	return result
}

func InitStmts() {
}