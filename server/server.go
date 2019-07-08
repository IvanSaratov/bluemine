package server

import (
	"html/template"
	"log"

	"github.com/IvanSaratov/bluemine/config"
	"github.com/gorilla/sessions"
	"github.com/jmoiron/sqlx"
)

//Core struct contains main vars of server
var Core struct {
	DB        *sqlx.DB
	Store     *sessions.CookieStore
	Templates map[string]*template.Template
}

//Init function initializes server
func Init() {
	Core.DB = sqlx.MustConnect("postgres", "host="+config.Conf.Host+" port="+config.Conf.DBPort+" user="+config.Conf.DBUser+" dbname="+config.Conf.DBName+" sslmode=disable password="+config.Conf.DBPassword)
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
