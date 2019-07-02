package server

import (
	"database/sql"
	"errors"
	"log"

	"github.com/gorilla/sessions"

	"github.com/IvanSaratov/bluemine/config"
)

//Core struct contains main vars of server
var Core struct {
	DB    *sql.DB
	Store *sessions.CookieStore
}

//Init function initializes server
func Init() (err error) {
	Core.DB, err = sql.Open("postgres", config.Conf.Postgresql)
	if err != nil {
		return err
	}
	log.Println("Connected to database")

	if config.Conf.SessionKey == "" {
		return errors.New("Empty session key")
	}
	Core.Store = sessions.NewCookieStore([]byte(config.Conf.SessionKey))
	log.Println("Created cookie store")

	return nil
}
