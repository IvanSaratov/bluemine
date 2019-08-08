package helpers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/IvanSaratov/bluemine/data"
	"github.com/IvanSaratov/bluemine/server"
)

//AlreadyLogin check user log status
func AlreadyLogin(r *http.Request) bool {
	session, _ := server.Core.Store.Get(r, "bluemine_session")
	err := server.Core.DB.QueryRow("SELECT id from profiles WHERE id = $1", session.Values["userid"])
	return err != nil
}

//AuthCheck is a middleware for handlers
func AuthCheck(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !AlreadyLogin(r) {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		handler(w, r)
	}
}

//ConvertIDToExecName convert executor's ID to executor's name
func ConvertIDToExecName(ID int, executorType string) (string, error) {
	var (
		executor string
		err      error
	)

	switch executorType {
	case "user":
		err = server.Core.DB.QueryRow("SELECT username FROM profiles WHERE id = $1", ID).Scan(&executor)
	case "group":
		err = server.Core.DB.QueryRow("SELECT group_name FROM groups WHERE id = $1", ID).Scan(&executor)
	default:
		return "", errors.New("Wrong executor_type")
	}

	return executor, err
}

//ConvertIDToExecFIO convert executor's ID to executor's fio
func ConvertIDToExecFIO(ID int) (string, error) {
	var (
		executor string
		err      error
	)

	err = server.Core.DB.QueryRow("SELECT user_fio FROM profiles WHERE id = $1", ID).Scan(&executor)

	return executor, err
}

//ConvertExecNameToID convert executor's name to executor's ID
func ConvertExecNameToID(executor string, executorType string) (int, error) {
	var (
		executorID int
		err        error
	)

	switch executorType {
	case "user":
		err = server.Core.DB.QueryRow("SELECT id FROM profiles WHERE username = $1", executor).Scan(&executorID)
	case "group":
		err = server.Core.DB.QueryRow("SELECT id FROM groups WHERE group_name = $1", executor).Scan(&executorID)
	default:
		return -1, errors.New("Wrong executor_type")
	}

	return executorID, err
}

//ConvertExecFIOToID convert executor's fio to executor's ID
func ConvertExecFIOToID(ID string) (int, error) {
	var (
		executorID int
		err        error
	)

	err = server.Core.DB.QueryRow("SELECT id FROM profiles WHERE user_fio = $1", ID).Scan(&executorID)

	return executorID, err
}

//GetCurrentUser gets info about current logged user
func GetCurrentUser(r *http.Request) (data.User, error) {
	session, _ := server.Core.Store.Get(r, "bluemine_session")

	var user data.User

	err := server.Core.DB.QueryRow("SELECT * FROM profiles WHERE username = $1", fmt.Sprint(session.Values["user"])).Scan(&user.UserID, &user.UserName, &user.UserFIO, &user.UserisAdmin, &user.UserRate)
	if err != nil {
		return user, err
	}

	return user, nil
}
