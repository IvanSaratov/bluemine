package main

import (
	"log"
	"net/http"

	"github.com/IvanSaratov/bluemine/config"
	"github.com/IvanSaratov/bluemine/handlers"
	"github.com/IvanSaratov/bluemine/server"

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
	router.HandleFunc("/makeadmin", handlers.MakeAdminHandler).Methods("POST")
	router.HandleFunc("/removeadmin", handlers.RemoveAdminHandler).Methods("POST")
	router.HandleFunc("/group/show/{id}", handlers.GroupHandler)
	router.HandleFunc("/group/new", handlers.AddGroupHandler).Methods("POST")
	router.HandleFunc("/groups", handlers.GroupsHandler)
	router.HandleFunc("/group/change", handlers.GroupChangeHandler)
	router.HandleFunc("/tasks", handlers.TasksHandler)
	router.HandleFunc("/tasks/show/{id}", handlers.TaskPageHandler)
	router.HandleFunc("/gettaskdata", handlers.GetTaskData).Methods("GET")
	router.HandleFunc("/gettaskdesc", handlers.GetTaskDesc).Methods("GET")
	router.HandleFunc("/tasks/new", handlers.AddTaskHandler).Methods("POST")
	router.HandleFunc("/tasks/change", handlers.ChangeTaskHandler).Methods("POST")
	router.HandleFunc("/tmpl/new", handlers.AddTmplHandler).Methods("POST")
	router.HandleFunc("/gettmpldata", handlers.GetTmplData).Methods("GET")
	router.HandleFunc("/tasks/close", handlers.TaskCloseHandler)
	router.HandleFunc("/tasks/open", handlers.TaskOpenHandler).Methods("POST")
	router.HandleFunc("/wiki", handlers.WikiHandler)
	router.HandleFunc("/wiki/show", handlers.AddWikiHandler)
	router.HandleFunc("/wiki/new", handlers.AddWikiHandler).Methods("POST")
	router.HandleFunc("/", handlers.RootHandler)

	log.Printf("Server listening on %s port", config.Conf.ListenPort)
	log.Fatal(http.ListenAndServe(config.Conf.ListenPort, router))
}
