package handlers

import (
	"html/template"
	"net/http"

	"github.com/IvanSaratov/bluemine/data"
)

func UserProfileHandler(w http.ResponseWriter, r *http.Request) {
	data := data.ViewData{
		UserData: data.User{
			UserName:       "test",
			UserFIO:        "test_testovich",
			UserDepartment: 2,
		},
	}

	tmpl, _ := template.ParseFiles("public/html/profile.html")
	tmpl.Execute(w, data)
}
