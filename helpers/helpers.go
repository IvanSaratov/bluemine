package helpers

import (
	"errors"
	"net/http"

	"github.com/IvanSaratov/bluemine/server"
)

//AlreadyLogin check user log status
func AlreadyLogin(r *http.Request) bool {
	session, _ := server.Core.Store.Get(r, "bluemine_session")
	return session.Values["userName"] != nil
}

func ConvertIDandStat(ID, stat int, executor_type string) (string, string, error) {
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
		return executor, "", errors.New("Wrong executor_type")
	}

	switch stat {
	case 0:
		return executor, "В процессе", err
	case 1:
		return executor, "Отложена", err
	case 2:
		return executor, "Закрыта", err
	default:
		return executor, "", errors.New("Wrong stat")
	}
}
