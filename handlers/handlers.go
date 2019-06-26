package handlers

import (
	"crypto/md5"
	"crypto/sha1"
	"errors"
	"fmt"
	"html/template"
	"io"
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

func loginUser(userLogin, userPassword string) (sessionId string, err error) {
	var password string

	if passwordHash(userPassword) != password {
		err = errors.New("Incorrect password")
		return
	}

	return
}

func passwordHash(password string) string {
	sh := sha1.New()
	io.WriteString(sh, password)

	md := md5.New()
	io.WriteString(md, password)

	return fmt.Sprintf("%x:%x", sh.Sum(nil), md.Sum(nil))
}
