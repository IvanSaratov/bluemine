package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"

	"github.com/IvanSaratov/bluemine/config"
	"github.com/IvanSaratov/bluemine/handlers"

	//"github.com/IvanSaratov/bluemine/session"

	_ "github.com/cockroachdb/cockroach-go/crdb"
	"github.com/gorilla/mux"
)

func main() {
	var configPath string

	flag.StringVar(&configPath, "c", "conf.toml", "Path to server configuration")
	flag.Parse()

	err := config.ParceConfig(configPath)
	if err != nil {
		log.Fatal("Error parsing config: ", err)
	}

	db, err := sql.Open("postgres", config.Conf.Postgresql)
	if err != nil {
		log.Fatal("Can't connect to database " + err.Error())
	}
	defer db.Close()

	router := mux.NewRouter()

	router.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("public"))))
	router.Handle("/login/{type}", handlers.AuthHandler)

	log.Fatal(http.ListenAndServe(config.Conf.Bind, router))

	var nilCh chan bool
	<-nilCh
}
