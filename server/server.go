package server

import (
	"database/sql"
	"errors"
	"html/template"
	"log"

	"github.com/gorilla/sessions"

	"github.com/IvanSaratov/bluemine/config"
)

//Core struct contains main vars of server
var Core struct {
	DB        *sql.DB
	Store     *sessions.CookieStore
	Templates map[string]*template.Template
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

	Core.Templates = make(map[string]*template.Template)
	temp := template.Must(template.ParseFiles("public/html/layout.html", "public/html/tasks.html"))
	Core.Templates["tasks"] = temp
	temp = template.Must(template.ParseFiles("public/html/layout.html", "public/html/profile.html"))
	Core.Templates["profile"] = temp
	temp = template.Must(template.ParseFiles("public/html/layout.html", "public/html/addtask.html"))
	Core.Templates["addTask"] = temp
	temp = template.Must(template.ParseFiles("public/html/layout.html", "public/html/taskpage.html"))
	Core.Templates["taskPage"] = temp
	log.Println("All templates parsed")

	return nil
}
