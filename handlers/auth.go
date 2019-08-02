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

func auth(login, password string) (int64, string, error) {
	if password == "" {
		return 0, "", errors.New("Empty password")
	}

	l, err := ldap.Dial("tcp", config.Conf.LdapServer)
	if err != nil {
		return 0, "", err
	}
	defer l.Close()

	if err := l.Bind(config.Conf.LdapUser, config.Conf.LdapPassword); err != nil {
		return 0, "", err
	}

	searchRequest := ldap.NewSearchRequest(config.Conf.LdapBaseDN, ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false, "(&(sAMAccountName="+login+"))", []string{"cn"}, nil)

	sr, err := l.Search(searchRequest)
	if err != nil {
		return 0, "", err
	}
	if len(sr.Entries) != 1 {
		return 0, "", errors.New("User not found")
	}

	username := sr.Entries[0].GetAttributeValue("cn")

	if err = l.Bind(username, password); err != nil {
		return 0, "", err
	}

	var userID int64

	err = userExists(login)
	if err != nil {
		if err != sql.ErrNoRows {
			return 0, "", err
		}

		userID, err = db.RegisterUser(server.Core.DB, l, login, username)
		if err != nil {
			return 0, "", err
		}
	} else {
		err = server.Core.DB.QueryRow("SELECT id FROM profiles WHERE username = $1", login).Scan(&userID)
		if err != nil {
			return 0, "", err
		}
	}

	return userID, username, err
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
		err := server.Core.Rnd.HTML(w, http.StatusOK, "login", nil)
		if err != nil {
			log.Print("Error parse template: ", err)
		}
	} else if r.Method == "POST" {
		r.ParseForm()
		login := r.FormValue("username")
		password := r.FormValue("password")

		login = strings.ToLower(login)

		userID, userName, err := auth(login, password)
		if err != nil {
			log.Printf("Failed to log in from %s using \"%s\" username: %s", r.RemoteAddr, login, err)

			http.Redirect(w, r, "/login", http.StatusMovedPermanently)
			return
		}

		log.Printf("User \"%s\" successfully log in from %s", login, r.RemoteAddr)
		session.Values["userid"] = userID
		session.Values["username"] = userName
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

	session.Values["userid"] = nil
	session.Values["username"] = nil
	session.Values["user"] = nil

	err = session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/login", http.StatusFound)
}
