package handlers

import (
	"html/template"
	"net/http"

	"github.com/IvanSaratov/bluemine/data"
	"github.com/IvanSaratov/bluemine/helpers"
	"github.com/IvanSaratov/bluemine/server"
)

//UserProfileHandler handle user's profile page
func UserProfileHandler(w http.ResponseWriter, r *http.Request) {
	if !helpers.AlreadyLogin(r) {
		http.Redirect(w, r, "/login", 301)
		return
	}

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
	if !helpers.AlreadyLogin(r) {
		http.Redirect(w, r, "/login", 301)
		return
	}

	data := data.ViewData{
		UserData: data.User{
			UserName:       "test",
			UserFIO:        "test_testovich",
			UserDepartment: "Otdel_Debilov",
		},
		TaskData: data.Task{
			TaskName:     "test",
			TaskExecutor: "Lox",
			TaskStat:     "V_Pizde",
		},
	}

	tmpl, _ := template.ParseFiles("public/html/tasks.html")
	tmpl.Execute(w, data)
}

//TaskPageHandler handle page of task
func TaskPageHandler(w http.ResponseWriter, r *http.Request) {
	if !helpers.AlreadyLogin(r) {
		http.Redirect(w, r, "/login", 301)
		return
	}

	data := data.ViewData{
		UserData: data.User{
			UserName:       "test",
			UserFIO:        "test_testovich",
			UserDepartment: "Otdel_Debilov",
		},
		TaskData: data.Task{
			TaskName:     "test",
			TaskDescPath: "/private/docs/test.txt",
			TaskExecutor: "Lox",
			TaskStat:     "V_Pizde",
		},
	}

	tmpl, _ := template.ParseFiles("public/html/taskpage.html")
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
		//TODO: Transform Status and ExecutorID from int to string
		/*var stat int
		var executerID int*/
		var (
			stat       string
			executorID string
		)

		err = rows.Scan(&name, &stat, &executorID)
		if err != nil {
			return err
		}

		tasksList = append(tasksList, data.Task{TaskName: name, TaskStat: stat, TaskExecutor: executorID})
	}

	return nil
}
