package handlers

import (
	"html/template"
	"net/http"

	"github.com/IvanSaratov/bluemine/data"
	"github.com/IvanSaratov/bluemine/server"
)

//UserProfileHandler handle user's profile page
func UserProfileHandler(w http.ResponseWriter, r *http.Request) {
	data := data.ViewData{
		UserData: data.User{
			UserName:       "test",
			UserFIO:        "test_testovich",
			UserDepartment: "Otdel_Debilov",
		},
	}

	tmpl, _ := template.ParseFiles("public/html/profile.html")
	tmpl.Execute(w, data)
}

//TasksHandler handle page with tasks
func TasksHandler(w http.ResponseWriter, r *http.Request) {
	data := data.ViewData{
		UserData: data.User{
			UserName:       "test",
			UserFIO:        "test_testovich",
			UserDepartment: "Otdel_Debilov",
		},
		TaskData: data.Task{
			TaskName: "test",
			TaskExecuter: 1337,
			TaskStat: 1,
		},
	}

	tmpl, _ := template.ParseFiles("public/html/tasks.html")
	tmpl.Execute(w, data)
}

func readTasks() error {
	rows, err := server.Core.DB.Query("SELECT task_name, stat, executor_id FROM tasks")
	if err != nil {
		return err
	}
	defer rows.Close()

	var tasksList []data.Task
	for rows.Next() {
		var name string
		var stat int
		var executerID int

		err = rows.Scan(&name, &stat, &executerID)
		if err != nil {
			return err
		}

		tasksList = append(tasksList, data.Task{TaskName: name, TaskStat: stat, TaskExecuter: executerID})
	}

	return nil
}
