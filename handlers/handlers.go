package handlers

import (
	"log"
	"net/http"
	"strconv"

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

	currentUser, err := helpers.GetCurrentUser(r)
	if err != nil {
		log.Printf("Error getting current user: %s", err)
	}

	vars := mux.Vars(r)
	username := vars["user"]

	user, err := db.GetUserInfo(server.Core.DB, username)
	if err != nil {
		log.Printf("Error getting info about %s: %s", username, err)
	}

	data := data.ViewData{
		CurrentUser: currentUser,
		UserData:    user,
	}

	err = server.Core.Templates["profile"].ExecuteTemplate(w, "base", data)
	if err != nil {
		log.Print(err)
	}
}

//TasksHandler handle page with tasks
func TasksHandler(w http.ResponseWriter, r *http.Request) {
	if !helpers.AlreadyLogin(r) {
		http.Redirect(w, r, "/login", 301)
		return
	}

	currentUser, err := helpers.GetCurrentUser(r)
	if err != nil {
		log.Printf("Error getting current user: %s", err)
	}

	tasks, err := db.GetAllTasks(server.Core.DB)
	if err != nil {
		log.Printf("Error getting tasks list: %s", err)
	}

	data := data.ViewData{
		CurrentUser: currentUser,
		Tasks:       tasks,
	}

	err = server.Core.Templates["tasks"].ExecuteTemplate(w, "base", data)
	if err != nil {
		log.Print(err)
	}
}

//TaskPageHandler handle page of task
func TaskPageHandler(w http.ResponseWriter, r *http.Request) {
	if !helpers.AlreadyLogin(r) {
		http.Redirect(w, r, "/login", 301)
		return
	}

	currentUser, err := helpers.GetCurrentUser(r)
	if err != nil {
		log.Printf("Error getting current user: %s", err)
	}

	vars := mux.Vars(r)
	taskIDstr := vars["id"]

	taskIDint, err := strconv.Atoi(taskIDstr)
	if err != nil {
		log.Printf("Error converting string to int on %s task page: %s", taskIDstr, err)
	}

	task, err := db.GetTask(server.Core.DB, taskIDint)
	if err != nil {
		log.Printf("Error getting task info from DB on %s task page: %s", taskIDstr, err)
	}

	data := data.ViewData{
		CurrentUser: currentUser,
		TaskData:    task,
	}

	err = server.Core.Templates["taskPage"].ExecuteTemplate(w, "base", data)
	if err != nil {
		log.Print(err)
	}
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
			executorID   int
			executorType string
			stat         string
			executor     string
		)

		err = rows.Scan(&taskID, &name, &stat, &executorID, &executorType)
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

	err := server.Core.DB.QueryRow("SELECT * FROM tasks WHERE id = $1", taskID).Scan(&task.TaskID, &task.TaskName, &executorType, &executorID,
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
