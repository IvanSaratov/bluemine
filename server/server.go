package server

import (
	"database/sql"
	"html/template"
	"log"

	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
)

//Core struct contains main vars of server
var Core struct {
	DB        *sql.DB
	Store     *sessions.CookieStore
	Templates map[string]*template.Template
}

//Init function initializes server
func init() {
	authKeyOne := securecookie.GenerateRandomKey(64)
	encryptionKeyOne := securecookie.GenerateRandomKey(32)

	Core.Store = sessions.NewCookieStore(
		authKeyOne,
		encryptionKeyOne,
	)

	Core.Store.Options = &sessions.Options{
		MaxAge:   24 * 60 * 60,
		HttpOnly: true,
	}

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
}
