package server

import (
	"database/sql"
	"errors"

	"github.com/gorilla/sessions"

	"github.com/IvanSaratov/bluemine/config"
)

var Core struct {
	DB    *sql.DB
	Store *sessions.CookieStore
}

func Init() (err error) {
	Core.DB, err = sql.Open("postgres", config.Conf.Postgresql)
	if err != nil {
		return err
	}

	if config.Conf.SessionKey == "" {
		return errors.New("Empty session key")
	}
	Core.Store = sessions.NewCookieStore([]byte(config.Conf.SessionKey))

	return nil
}
