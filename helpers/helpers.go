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

func ConvertIDToExec(ID int, executor_type string) (string, error) {
	var (
		executor string
		err      error
	)

	switch executor_type {
	case "user":
		err = server.Core.DB.QueryRow("SELECT user_fio FROM profiles WHERE id = $1", ID).Scan(&executor)
	case "group":
		err = server.Core.DB.QueryRow("SELECT group_name FROM groups WHERE id = $1", ID).Scan(&executor)
	default:
		return "", errors.New("Wrong executor_type")
	}

	return executor, err
}

func ConvertExecToID(executor string, executor_type string) (int, error) {
	var (
		executorID int
		err        error
	)

	switch executor_type {
	case "user":
		err = server.Core.DB.QueryRow("SELECT id FROM profiles WHERE user_fio = $1", executor).Scan(&executorID)
	case "group":
		err = server.Core.DB.QueryRow("SELECT id FROM groups WHERE group_name = $1", executor).Scan(&executorID)
	default:
		return -1, errors.New("Wrong executor_type")
	}

	return executorID, err
}

func GetCurrentUser(r *http.Request) data.User {
	session, _ := server.Core.Store.Get(r, "bluemine_session")

	return data.User{UserName: fmt.Sprint(session.Values["user"]), UserFIO: fmt.Sprint(session.Values["username"])}
}
