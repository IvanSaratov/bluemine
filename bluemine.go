package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/IvanSaratov/bluemine/config"
	"github.com/IvanSaratov/bluemine/handlers"
	"github.com/IvanSaratov/bluemine/helpers"
	"github.com/IvanSaratov/bluemine/server"

	"github.com/braintree/manners"
	_ "github.com/cockroachdb/cockroach-go/crdb"
	"github.com/gorilla/mux"
)

func main() {
	var configPath string

	log.Println("Starting...")

	err := config.ParseConfig(configPath)
	if err != nil {
		log.Fatal("Error parsing config: ", err)
	}
	log.Println("Config parsed!")

	defer server.Core.DB.Close()

	router := mux.NewRouter()

	router.PathPrefix("/private/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !helpers.AlreadyLogin(r) {
			http.Redirect(w, r, "/login", 302)
			return
		}

		realHandler := http.StripPrefix("/private/", http.FileServer(http.Dir("./private/"))).ServeHTTP
		realHandler(w, r)
	})
	router.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("./public/"))))
	router.HandleFunc("/login", handlers.LoginHandler)
	router.HandleFunc("/logout", handlers.LogoutHandler)
	router.HandleFunc("/profile/{user}", handlers.UserProfileHandler)
	router.HandleFunc("/group/{group}", handlers.GroupHandler)
	router.HandleFunc("/tasks", handlers.TasksHandler)
	router.HandleFunc("/tasks/show/{id}", handlers.TaskPageHandler)
	router.HandleFunc("/tasks/new", handlers.AddTaskHandler)
	router.HandleFunc("/tasks/close", handlers.TaskCloseHandler)
	router.HandleFunc("/wiki/new", handlers.AddWikiHandler)
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if !helpers.AlreadyLogin(r) {
			http.Redirect(w, r, "/login", 301)
		} else {
			session, _ := server.Core.Store.Get(r, "bluemine_session")
			http.Redirect(w, r, "/profile/"+fmt.Sprintf("%v", session.Values["user"]), 301)
		}
	})

	ch := make(chan os.Signal)
	signal.Notify(ch, os.Interrupt, os.Kill)
	go func(ch <-chan os.Signal) {
		<-ch
		manners.Close()
	}(ch)

	log.Printf("Server listening on %s port", config.Conf.Bind)
	log.Fatal(manners.ListenAndServe(config.Conf.Bind, router))
}
