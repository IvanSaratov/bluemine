package main

import (
	"log"
	"net/http"

	"github.com/IvanSaratov/bluemine/config"
	"github.com/IvanSaratov/bluemine/handlers"
	"github.com/IvanSaratov/bluemine/server"

	//"golang.org/x/net/http2"
	_ "github.com/cockroachdb/cockroach-go/crdb"
	"github.com/gorilla/mux"
)

func init() {
	var configPath = "conf.toml"

	log.Println("Starting...")

	err := config.ParseConfig(configPath)
	if err != nil {
		log.Fatal("Error parsing config: ", err)
	}
	log.Println("Config parsed!")

	server.Init()
}

func main() {
	defer server.Core.DB.Close()

	router := mux.NewRouter()

	router.PathPrefix("/private/").HandlerFunc(handlers.PrivateHandler)
	router.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("./public/"))))
	router.HandleFunc("/login", handlers.LoginHandler)
	router.HandleFunc("/logout", handlers.LogoutHandler)
	router.HandleFunc("/profile/{user}", handlers.UserProfileHandler)
	router.HandleFunc("/group/show/{id}", handlers.GroupHandler)
	router.HandleFunc("/group/new", handlers.AddGroupHandler)
	router.HandleFunc("/groups", handlers.GroupsHandler)
	router.HandleFunc("/tasks", handlers.TasksHandler)
	router.HandleFunc("/tasks/show/{id}", handlers.TaskPageHandler)
	router.HandleFunc("/gettaskdesc", handlers.GetTaskDesc)
	router.HandleFunc("/tasks/new", handlers.AddTaskHandler).Methods("POST")
	router.HandleFunc("/tmpl/new", handlers.AddTmplHandler).Methods("POST")
	router.HandleFunc("/gettmpldata", handlers.GetTmplData).Methods("GET")
	router.HandleFunc("/tasks/close", handlers.TaskCloseHandler)
	//router.HandleFunc("/wiki/new", handlers.AddWikiHandler)
	router.HandleFunc("/", handlers.RootHandler)

	log.Printf("Server listening on %s port", config.Conf.ListenPort)
	log.Fatal(http.ListenAndServe(config.Conf.ListenPort, router))
}
