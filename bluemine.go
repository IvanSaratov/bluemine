package main

import (
	"crypto/md5"
	"crypto/sha1"
	"database/sql"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/IvanSaratov/bluemine/config"
	"github.com/IvanSaratov/bluemine/session"
	_ "github.com/cockroachdb/cockroach-go/crdb"
)

func main() {
	var (
		err        error
		configPath string
	)

	flag.StringVar(&configPath, "c", "conf.toml", "Path to server configuration")
	flag.Parse()

	config.ParceConfig(configPath)

	db, err := sql.Open("postgres", config.Conf.Postgresql)
	if err != nil {
		log.Fatal("Can't connect to database " + err.Error())
	}
	defer db.Close()

	http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("public"))))

	http.HandleFunc("/login", LoginHandler)
	http.HandleFunc("/logout", LogoutHandler)

	go listen(config.Conf.Bind)

	var nilCh chan bool
	<-nilCh
}

func listen(addr string) {
	log.Fatal("ListenAndServe: ", http.ListenAndServe(addr, nil))
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
