package handlers

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/IvanSaratov/bluemine/data"
	"github.com/IvanSaratov/bluemine/db"
	"github.com/IvanSaratov/bluemine/helpers"
	"github.com/IvanSaratov/bluemine/server"
	"github.com/gorilla/mux"
)

//UserProfileHandler handle user's profile page
func UserProfileHandler(w http.ResponseWriter, r *http.Request) {
	if !helpers.AlreadyLogin(r) {
		http.Redirect(w, r, "/login", 301)
		return
	}
	session, _ := server.Core.Store.Get(r, "bluemine_session")

	vars := mux.Vars(r)
	username := vars["user"]

	user, err := db.GetUserInfo(server.Core.DB, username)
	if err != nil {
		log.Printf("Error getting info about %s: %s", username, err)
	}

	data := data.ViewData{
		CurrentUser: data.User{
			UserName: fmt.Sprint(session.Values["user"]),
			UserFIO:  fmt.Sprint(session.Values["username"]),
		},
		UserData: user,
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

	taskList, err := readTasks()
	if err != nil {
		log.Printf("Error reading tasks: %s", err)
	}

	data := data.ViewData{
		UserData: data.User{
			UserName:       "test",
			UserFIO:        "test_testovich",
			UserDepartment: "Otdel_Debilov",
		},
		Tasks: taskList,
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
			TaskExecutor: "Lox",
			TaskStat:     "In progress",
		},
	}

	tmpl, _ := template.ParseFiles("public/html/taskpage.html")
	tmpl.Execute(w, data)
}

func readTasks() ([]data.Task, error) {
	rows, err := server.Core.DB.Query("SELECT id, task_name, stat, executor_id, executor_type FROM tasks")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasksList []data.Task
	for rows.Next() {
		var (
			taskID       int
			name         string
			statInt      int
			executorID   int
			executorType string
			stat         string
			executor     string
		)

		err = rows.Scan(&taskID, &name, &statInt, &executorID, &executorType)
		if err != nil {
			return nil, err
		}

		executor, err = helpers.ConvertIDToExec(executorID, executorType)
		if err != nil {
			return nil, err
		}

		tasksList = append(tasksList, data.Task{TaskID: taskID, TaskName: name, TaskStat: stat, TaskExecutor: executor})
	}

	return tasksList, nil
}

func readTask(taskID int) (data.Task, error) {
	var (
		task         data.Task
		executorID   int
		executorType string
	)

	err := server.Core.DB.QueryRow("SELECT * FROM tasks WHERE id = $1", taskID).Scan(&task.TaskID, &task.TaskName, &task.TaskDescPath, &executorType, &executorID,
		&task.TaskStat, &task.TaskDateStart, &task.TaskDateEnd, &task.TaskRate)
	if err != nil {
		return task, err
	}

	task.TaskExecutor, err = helpers.ConvertIDToExec(executorID, executorType)
	if err != nil {
		return task, err
	}

	return task, nil
}
