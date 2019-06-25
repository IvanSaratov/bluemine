package handler

import (
	"crypto/md5"
	"crypto/sha1"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

func LoginHandler(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	userLogin := req.Form.Get("login")
	userPassword := req.Form.Get("password")

	if userLogin == "" || userPassword == "" {
		fmt.Fprintf(w, "Empty login or password")
		return
	}

	sessionId, err := loginUser(userLogin, userPassword)
	if err != nil {
		return
	}

	cookie := &http.Cookie{
		Name:    "id",
		Value:   string(sessionId),
		Path:    "/",
		Domain:  req.Header.Get("Host"),
		Expires: time.Now().Add(360 * 24 * time.Hour),
	}

	http.SetCookie(w, cookie)
	w.Header().Add("Location", "/")
	w.WriteHeader(302)
}

func loginUser(userLogin, userPassword string) (sessionId string, err error) {
	var password string

	if passwordHash(userPassword) != password {
		err = errors.New("Incorrect password")
		return
	}

	return
}

func LogoutHandler(w http.ResponseWriter, req *http.Request) {
	http.SetCookie(w, &http.Cookie{Name: "id"})
	w.Header().Add("Location", "/")
	w.WriteHeader(302)
}

func passwordHash(password string) string {
	sh := sha1.New()
	io.WriteString(sh, password)

	md := md5.New()
	io.WriteString(md, password)

	return fmt.Sprintf("%x:%x", sh.Sum(nil), md.Sum(nil))
}
