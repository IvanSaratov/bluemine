package server

import (
	"log"

	"github.com/IvanSaratov/bluemine/config"
	"github.com/gorilla/sessions"
	"github.com/jmoiron/sqlx"
	"github.com/thedevsaddam/renderer"
)

//Core struct contains main vars of server
var Core struct {
	DB    *sqlx.DB
	Store *sessions.CookieStore
	Rnd   *renderer.Render
}

//Init function initializes server
func Init() {
	Core.DB = sqlx.MustConnect("postgres", "user="+config.Conf.DBUser+" password="+config.Conf.DBPassword+" host="+config.Conf.Host+" port="+config.Conf.DBPort+" dbname="+config.Conf.DBName)
	log.Println("Connected to database successfull")

	Core.Store = sessions.NewCookieStore(
		[]byte(config.Conf.SessionKey),
		[]byte(config.Conf.EncryptionKey),
	)

	Core.Store.Options = &sessions.Options{
		HttpOnly: true,
	}
	log.Println("Created cookie store")

	Core.Rnd = renderer.New(
		renderer.Options{
			ParseGlobPattern: "html/*.html",
		},
	)

	log.Println("All templates parsed")
}
