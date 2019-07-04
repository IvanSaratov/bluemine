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
