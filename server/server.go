package server

import (
	"database/sql"
	"html/template"
	"log"

	"github.com/IvanSaratov/bluemine/config"
	"github.com/gorilla/sessions"
)

//Core struct contains main vars of server
var Core struct {
	DB        *sql.DB
	Store     *sessions.CookieStore
	Templates map[string]*template.Template
}

//Init function initializes server
func Init() {
	var err error

	Core.DB, err = sql.Open("postgres", config.Conf.DBHost)
	if err != nil {
		log.Fatal("Error connecting to database: ", err)
	}
	log.Println("Connected to database successfull")

	Core.Store = sessions.NewCookieStore(
		[]byte(config.Conf.SessionKey),
		[]byte(config.Conf.EncryptionKey),
	)

	Core.Store.Options = &sessions.Options{
		MaxAge:   24 * 60 * 60,
		HttpOnly: true,
	}
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
	temp = template.Must(template.ParseFiles("public/html/layout.html", "public/html/group.html"))
	Core.Templates["group"] = temp

	log.Println("All templates parsed")
}
