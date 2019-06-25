package main

import (
	"database/sql"
	"flag"
	"log"
	"os"
	"net/http"
	"io"
	"strings"

	"github.com/IvanSaratov/bluemine/config"
	"github.com/IvanSaratov/bluemine/handlers"

	//"github.com/IvanSaratov/bluemine/session"

	"github.com/gorilla/mux"
	_ "github.com/cockroachdb/cockroach-go/crdb"
)

func main() {
	var (
		configPath string
	)

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

	router.HandleFunc("/login", handlers.LoginHandler)
	router.HandleFunc("/logout", handlers.LoginHandler)

	go listen(config.Conf.Bind)

	var nilCh chan bool
	<-nilCh
}

func listen(addr string) {
	log.Fatal("ListenAndServe: ", http.ListenAndServe(addr, nil))
}

func serveStatic(filename string, w http.ResponseWriter) {
	file, err := os.Open("filename")
	if err != nil {
		w.WriteHeader(404)
		w.Write([]byte("Could not find file: " + filename))
		return
	}
	defer file.Close()

	if strings.HasSuffix(filename, ".css") {
		w.Header().Add("Content-type", "text/css")
	} else if strings.HasSuffix(filename, ".js") {
		w.Header().Add("Content-type", "application/javascript")
	}

	io.Copy(w, file)

}