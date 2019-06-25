package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"

	"github.com/IvanSaratov/bluemine/config"
	"github.com/IvanSaratov/bluemine/handler"
	//"github.com/IvanSaratov/bluemine/session"
	_ "github.com/cockroachdb/cockroach-go/crdb"
)

func main() {
	var (
		err        error
		configPath string
	)

	flag.StringVar(&configPath, "c", "conf.toml", "Path to server configuration")
	flag.Parse()

	config.ParceConfig(configPath)

	db, err := sql.Open("postgres", config.Conf.Postgresql)
	if err != nil {
		log.Fatal("Can't connect to database " + err.Error())
	}
	defer db.Close()

	http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("public"))))

	http.HandleFunc("/login", handler.LoginHandler)
	http.HandleFunc("/logout", handler.LoginHandler)

	go listen(config.Conf.Bind)

	var nilCh chan bool
	<-nilCh
}

func listen(addr string) {
	log.Fatal("ListenAndServe: ", http.ListenAndServe(addr, nil))
}
