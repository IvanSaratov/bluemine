package handlers

import (
	"html/template"
	"net/http"

	"github.com/IvanSaratov/bluemine/data"
	"github.com/IvanSaratov/bluemine/helpers"
)

//AddTaskHandler handle task adding page
func AddTaskHandler(w http.ResponseWriter, r *http.Request) {
	if !helpers.AlreadyLogin(r) {
		http.Redirect(w, r, "/login", 301)
		return
	}

	if r.Method == "GET" {
		data := data.ViewData{
			UserData: data.User{
				UserName:       "test",
				UserFIO:        "test_testovich",
				UserDepartment: "Otdel_Debilov",
			},
		}

		tmpl, _ := template.ParseFiles("public/html/addtask.html")
		tmpl.Execute(w, data)
	}
}
