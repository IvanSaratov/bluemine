package main

import (
	"github.com/IvanSaratov/bluemine/config"
	"crypto/md5"
	"crypto/sha1"
	"database/sql"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	_ "github.com/cockroachdb/cockroach-go/crdb"
)

func main() {
	var (
		err        error
		configPath string
	)

	flag.StringVar(&configPath, "c", "config.toml", "Path to server configuration")
	flag.Parse()

	config.ParceConfig(configPath)

	db, err := sql.Open("postgres", config.Conf.Postgresql)
	if err != nil {
		log.Fatal("Can't connect to database " + err.Error())
	}
	defer db.Close()

	http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("public"))))

	http.Handle("/login", LoginHandler)
}

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
	var password, login string

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
