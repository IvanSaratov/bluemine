package handlers

import (
	"fmt"
	"html/template"
	"log"
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
	session, _ := server.Core.Store.Get(r, "bluemine_session")

	data := data.ViewData{
		UserData: data.User{
			UserName:       fmt.Sprint(session.Values["user"]),
			UserFIO:        fmt.Sprint(session.Values["username"]),
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

	tmpl, _ := template.ParseFiles("public/html/taskpage.html")
	tmpl.Execute(w, data)
}

func readTasks() ([]data.Task, error) {
	rows, err := server.Core.DB.Query("SELECT id, task_name, stat, executor_id FROM tasks")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasksList []data.Task
	for rows.Next() {
		var (
			taskID     int
			name       string
			statInt    int
			executorID int
			stat       string
			executor   string
		)

		err = rows.Scan(&taskID, &name, &statInt, &executorID)
		if err != nil {
			return nil, err
		}

		executor, stat, err = helpers.ConvertIDandStat(executorID, statInt)
		if err != nil {
			return nil, err
		}

		tasksList = append(tasksList, data.Task{TaskID: taskID, TaskName: name, TaskStat: stat, TaskExecutor: executor})
	}

	return tasksList, nil
}

func readTask(taskID int) (data.Task, error) {
	var (
		task       data.Task
		executorID int
		statInt    int
	)

	err := server.Core.DB.QueryRow("SELECT * FROM tasks WHERE id = $1", taskID).Scan(&task.TaskID, &task.TaskName, &task.TaskDescPath,
		&statInt, &task.TaskDateStart, &task.TaskDateEnd, &task.TaskRate, &executorID)
	if err != nil {
		return task, err
	}

	task.TaskExecutor, task.TaskStat, err = helpers.ConvertIDandStat(executorID, statInt)
	if err != nil {
		return task, err
	}

	return task, nil
}
