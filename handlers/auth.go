package handlers

import (
	"database/sql"
	"errors"
	"net/http"
	"strings"

	"github.com/IvanSaratov/bluemine/config"
	"github.com/IvanSaratov/bluemine/server"

	"github.com/go-ldap/ldap"
)

//AlreadyLogin check user log status
func AlreadyLogin(r *http.Request) bool {
	session, _ := server.Core.Store.Get(r, "bluemine_session")
	return session.Values["userName"] != nil
}

func auth(login, password string) (string, error) {
	if password == "" {
		return "", errors.New("Empty password")
	}

	l, err := ldap.Dial("tcp", config.Conf.LdapServer)
	if err != nil {
		return "", err
	}
	defer l.Close()

	if err := l.Bind(config.Conf.LdapUser, config.Conf.LdapPassword); err != nil {
		return "", err
	}

	searchRequest := ldap.NewSearchRequest(
		config.Conf.LdapBaseDN, ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false, "(&(sAMAccountName="+login+"))", []string{"cn"}, nil,
	)

	sr, err := l.Search(searchRequest)
	if err != nil {
		return "", err
	}
	if len(sr.Entries) != 1 {
		return "", errors.New("User not found")
	}

	username := sr.Entries[0].GetAttributeValue("cn")

	if err = l.Bind(username, password); err != nil {
		return "", err
	}

	login = strings.ToLower(login)
	err = userExists(login)
	if err != nil {
		if err != sql.ErrNoRows {
			return "", err
		}

		err = addUserToBD(login)
		if err != nil {
			return "", err
		}
	}

	return username, err
}

func addUserToBD(login string) error {
	createStmt := "INSERT INTO profiles (username, user_fio) VALUES ($1, $2)"
	_, err := server.Core.DB.Exec(createStmt, login, login+"_test")
	return err
}

func userExists(login string) error {
	userStmt := "SELECT username FROM profiles WHERE username = $1"
	err := server.Core.DB.QueryRow(userStmt, login).Scan(&login)
	return err
}

//LoginHandler handle login page
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := server.Core.Store.Get(r, "bluemine_session")

	if r.Method == "GET" {
		http.ServeFile(w, r, "public/html/login.html")
	} else if r.Method == "POST" {
		r.ParseForm()
		login := r.FormValue("username")
		password := r.FormValue("password")

		if AlreadyLogin(r) {
			http.Redirect(w, r, "/profile/"+login, 301)
			return
		}

		if userName, err := auth(login, password); err != nil {
			http.Redirect(w, r, "/login", 301)
			return
		} else {
			session.Values["userName"] = userName
			session.Values["user"] = login
			session.Save(r, w)
			http.Redirect(w, r, "/profile/"+login, 301)
		}
	}
}

//LogoutHandler handle logout page
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := server.Core.Store.Get(r, "bluemine_session")

	session.Values["userName"] = nil
	session.Values["user"] = nil
	session.Save(r, w)
	http.Redirect(w, r, "/login", 301)
}