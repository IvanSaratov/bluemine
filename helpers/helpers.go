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
	return session.Values["userName"] != nil
}

//ConvertIDToExec convert executor ID to executor string
func ConvertIDToExec(ID int, executorType string) (string, error) {
	var (
		executor string
		err      error
	)

	switch executorType {
	case "user":
		err = server.Core.DB.QueryRow("SELECT user_fio FROM profiles WHERE id = $1", ID).Scan(&executor)
	case "group":
		err = server.Core.DB.QueryRow("SELECT group_name FROM groups WHERE id = $1", ID).Scan(&executor)
	default:
		return "", errors.New("Wrong executor_type")
	}

	return executor, err
}

//ConvertExecToID convert executor string to executor ID
func ConvertExecToID(executor string, executorType string) (int, error) {
	var (
		executorID int
		err        error
	)

	switch executorType {
	case "user":
		err = server.Core.DB.QueryRow("SELECT id FROM profiles WHERE user_fio = $1", executor).Scan(&executorID)
	case "group":
		err = server.Core.DB.QueryRow("SELECT id FROM groups WHERE group_name = $1", executor).Scan(&executorID)
	default:
		return -1, errors.New("Wrong executor_type")
	}

	return executorID, err
}

//GetCurrentUser gets info about current logged user
func GetCurrentUser(r *http.Request) data.User {
	session, _ := server.Core.Store.Get(r, "bluemine_session")

	return data.User{UserName: fmt.Sprint(session.Values["user"]), UserFIO: fmt.Sprint(session.Values["username"])}
}
