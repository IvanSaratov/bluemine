package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/IvanSaratov/bluemine/config"
	"github.com/IvanSaratov/bluemine/handlers"
	"github.com/IvanSaratov/bluemine/server"

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

	err = server.Init()
	if err != nil {
		log.Fatal("Error initializing server: ", err)
	}
	defer server.Core.DB.Close()

	router := mux.NewRouter()

	router.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("./public/"))))
	router.HandleFunc("/login", handlers.LoginHandler)
	router.HandleFunc("/logout", handlers.LogoutHandler)

	log.Fatal(http.ListenAndServe(config.Conf.Bind, router))

	var nilCh chan bool
	<-nilCh
}
