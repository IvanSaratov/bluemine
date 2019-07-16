package handlers

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/IvanSaratov/bluemine/config"
	"github.com/IvanSaratov/bluemine/db"
	"github.com/IvanSaratov/bluemine/server"

	"github.com/go-ldap/ldap"
)

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

	searchRequest := ldap.NewSearchRequest(config.Conf.LdapBaseDN, ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false, "(&(sAMAccountName="+login+"))", []string{"cn"}, nil)

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

	err = userExists(login)
	if err != nil {
		if err != sql.ErrNoRows {
			return "", err
		}

		err = db.RegisterUser(server.Core.DB, l, login, username)
		if err != nil {
			return "", err
		}
	}

	return username, err
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
		err := server.Core.Templates["login"].ExecuteTemplate(w, "base", nil)
		if err != nil {
			log.Print("Error parse template: ", err)
		}
	} else if r.Method == "POST" {
		r.ParseForm()
		login := r.FormValue("username")
		password := r.FormValue("password")

		login = strings.ToLower(login)

		userName, err := auth(login, password)
		if err != nil {
			log.Printf("Failed to log in from %s using \"%s\" username: %s", r.RemoteAddr, login, err)

			http.Redirect(w, r, "/login", http.StatusMovedPermanently)
			return
		}

		log.Printf("User \"%s\" successfully log in from %s", login, r.RemoteAddr)
		session.Values["userName"] = userName
		session.Values["user"] = login
		session.Save(r, w)
		http.Redirect(w, r, "/profile/"+login, http.StatusMovedPermanently)
	}
}

//LogoutHandler handle logout page
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	session, err := server.Core.Store.Get(r, "bluemine_session")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	session.Values["userName"] = nil
	session.Values["user"] = nil
	session.Options.MaxAge = -1

	err = session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/login", http.StatusFound)
}
