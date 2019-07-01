package helpers

import (
	"database/sql"
	"net/http"

	"github.com/IvanSaratov/bluemine/server"
)

//AlreadyLogin check user log status
func AlreadyLogin(r *http.Request) bool {
	session, _ := server.Core.Store.Get(r, "bluemine_session")
	return session.Values["userName"] != nil
}

func ConvertIDandStat(ID, stat int) (string, string, error) {
	var executer string

	err := server.Core.DB.QueryRow("SELECT user_fio FROM profiles WHERE id = $1", ID).Scan(&executer)
	if err != nil {
		if err != sql.ErrNoRows {
			return "", "", err
		}
		err = server.Core.DB.QueryRow("SELECT group_name FROM groups WHERE id = $1", ID).Scan(&executer)
		if err != sql.ErrNoRows {
			return "", "", err
		}
	}

	switch stat {
	case 0:
		return executer, "В процессе", nil
	case 1:
		return executer, "Отложена", nil
	case 2:
		return executer, "Закрыта", nil
	default:
		return "", "", nil
	}
}
