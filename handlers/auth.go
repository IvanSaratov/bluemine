package handlers

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/IvanSaratov/bluemine/config"
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

		err = addUserToDB(login, l)
		if err != nil {
			return "", err
		}
	}

	return username, err
}

func addUserToDB(login string, l *ldap.Conn) error {
	result, err := l.Search(ldap.NewSearchRequest(
		config.Conf.LdapBaseDN,
		ldap.ScopeWholeSubtree,
		ldap.NeverDerefAliases,
		0,
		0,
		false,
		"(&(sAMAccountName="+login+"))",
		[]string{"memberOf", "cn"},
		nil,
	))
	if err != nil {
		return err
	}

	userFIO := result.Entries[0].GetAttributeValue("cn")
	var listOfMembers []string
	for _, x := range result.Entries[0].GetAttributeValues("memberOf") {
		listOfMembers = append(listOfMembers, x[strings.Index(x, "CN=")+3:strings.Index(x, ",")])
	}

	createUserStmt := "INSERT INTO profiles (id, username, user_fio) VALUES (DEFAULT, $1, $2) RETURNING id"
	var userID int64
	err = server.Core.DB.QueryRow(createUserStmt, login, userFIO).Scan(&userID)
	if err != nil {
		return err
	}

	rows, err := server.Core.DB.Query("SELECT id, group_name FROM groups")
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var groupID int
		var groupName string
		if err = rows.Scan(&groupID, &groupName); err != nil {
			log.Fatal(err)
		}
		for _, ldapGroupName := range listOfMembers {
			if groupName == ldapGroupName {
				_, err = server.Core.DB.Exec("INSERT INTO groups_profiles (group_id, profile_id) VALUES ($1, $2)", groupID, userID)
				if err != nil {
					return err
				}
				break
			}
		}
	}

	return nil
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

		userName, err := auth(login, password)
		if err != nil {
			log.Printf("Failed to log in from %s using \"%s\" username: %s", r.RemoteAddr, login, err)

			http.Redirect(w, r, "/login", 301)
			return
		}

		log.Printf("User \"%s\" successfully log in from %s", login, r.RemoteAddr)
		session.Values["userName"] = userName
		session.Values["user"] = login
		session.Save(r, w)
		http.Redirect(w, r, "/profile/"+login, 301)
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
