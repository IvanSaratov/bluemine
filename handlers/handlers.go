package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/IvanSaratov/bluemine/data"
	"github.com/IvanSaratov/bluemine/db"
	"github.com/IvanSaratov/bluemine/helpers"
	"github.com/IvanSaratov/bluemine/server"
	"github.com/gorilla/mux"
)

//RootHandler handle root path
func RootHandler(w http.ResponseWriter, r *http.Request) {
	if !helpers.AlreadyLogin(r) {
		http.Redirect(w, r, "/login", 301)
	} else {
		session, _ := server.Core.Store.Get(r, "bluemine_session")
		http.Redirect(w, r, "/profile/"+fmt.Sprintf("%v", session.Values["user"]), 301)
	}
}

//PrivateHandler handle private file server
func PrivateHandler(w http.ResponseWriter, r *http.Request) {
	if !helpers.AlreadyLogin(r) {
		http.Redirect(w, r, "/login", 302)
		return
	}

	realHandler := http.StripPrefix("/private/", http.FileServer(http.Dir("./private/"))).ServeHTTP
	realHandler(w, r)
}

//UserProfileHandler handle user's profile page
func UserProfileHandler(w http.ResponseWriter, r *http.Request) {
	if !helpers.AlreadyLogin(r) {
		http.Redirect(w, r, "/login", http.StatusMovedPermanently)
		return
	}

	currentUser, err := helpers.GetCurrentUser(r)
	if err != nil {
		log.Printf("Error getting current user: %s", err)
	}

	vars := mux.Vars(r)
	username := vars["user"]

	var userID int

	err = server.Core.DB.QueryRow("SELECT id FROM profiles WHERE username = $1", username).Scan(&userID)
	if err != nil {
		log.Printf("Error getting %s's id: %s", username, err)
	}

	user, err := db.GetUserbyID(server.Core.DB, userID)
	if err != nil {
		log.Printf("Error getting info about %s: %s", username, err)
	}

	viewData := data.ViewData{
		CurrentUser: currentUser,
		UserData:    user,
	}

	err = server.Core.Templates["profile"].ExecuteTemplate(w, "base", viewData)
	if err != nil {
		log.Print(err)
	}
}

//GroupHandler handle group's profile page
func GroupHandler(w http.ResponseWriter, r *http.Request) {
	if !helpers.AlreadyLogin(r) {
		http.Redirect(w, r, "/login", http.StatusMovedPermanently)
		return
	}

	currentUser, err := helpers.GetCurrentUser(r)
	if err != nil {
		log.Printf("Error getting current user: %s", err)
	}

	vars := mux.Vars(r)
	groupName := vars["group"]

	users, err := db.GetGroupUsers(server.Core.DB, groupName)
	if err != nil {
		log.Printf("Error getting info about %s: %s", groupName, err)
	}

	viewData := data.ViewData{
		CurrentUser: currentUser,
		Users:       users,
	}

	err = server.Core.Templates["group"].ExecuteTemplate(w, "base", viewData)
	if err != nil {
		log.Print(err)
	}
}

//TasksHandler handle page with tasks
func TasksHandler(w http.ResponseWriter, r *http.Request) {
	if !helpers.AlreadyLogin(r) {
		http.Redirect(w, r, "/login", http.StatusMovedPermanently)
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

	viewData := data.ViewData{
		CurrentUser: currentUser,
		Tasks:       tasks,
	}

	err = server.Core.Templates["tasks"].ExecuteTemplate(w, "base", viewData)
	if err != nil {
		log.Print(err)
	}
}

//TaskPageHandler handle page of task
func TaskPageHandler(w http.ResponseWriter, r *http.Request) {
	if !helpers.AlreadyLogin(r) {
		http.Redirect(w, r, "/login", http.StatusMovedPermanently)
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

	task, err := db.GetTaskbyID(server.Core.DB, taskIDint)
	if err != nil {
		log.Printf("Error getting task info from DB on %s task page: %s", taskIDstr, err)
	}

	viewData := data.ViewData{
		CurrentUser: currentUser,
		TaskData:    task,
	}

	err = server.Core.Templates["taskPage"].ExecuteTemplate(w, "base", viewData)
	if err != nil {
		log.Print(err)
	}
}

//TaskCloseHandler handle page to close task
func TaskCloseHandler(w http.ResponseWriter, r *http.Request) {
	if !helpers.AlreadyLogin(r) {
		http.Redirect(w, r, "/login", http.StatusMovedPermanently)
		return
	}

	taskID, _ := strconv.Atoi(r.FormValue("id"))
	task, err := db.GetTaskbyID(server.Core.DB, taskID)
	if err != nil {
		log.Print(err)
	}

	_, _ = server.Core.DB.Exec("UPDATE profiles SET rating = (rating + $1) WHERE user_fio = $2", task.TaskRate, task.TaskExecutorFIO)
	_, _ = server.Core.DB.Exec("UPDATE tasks SET stat = 'Закрыта' WHERE id = $1", task.TaskID)
}
