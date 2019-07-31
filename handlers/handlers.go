package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/IvanSaratov/bluemine/db"
	"github.com/IvanSaratov/bluemine/helpers"
	"github.com/IvanSaratov/bluemine/server"
	"github.com/gorilla/mux"
)

//RootHandler handle root path
func RootHandler(w http.ResponseWriter, r *http.Request) {
	if !helpers.AlreadyLogin(r) {
		http.Redirect(w, r, "/login", http.StatusMovedPermanently)
	} else {
		session, _ := server.Core.Store.Get(r, "bluemine_session")
		http.Redirect(w, r, "/profile/"+fmt.Sprintf("%v", session.Values["user"]), http.StatusMovedPermanently)
	}
}

//PrivateHandler handle private file server
func PrivateHandler(w http.ResponseWriter, r *http.Request) {
	if !helpers.AlreadyLogin(r) {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	realHandler := http.StripPrefix("/private/", http.FileServer(http.Dir("./private/"))).ServeHTTP
	realHandler(w, r)
}

//AdminActHandler handle user administrator status change
func AdminActHandler(w http.ResponseWriter, r *http.Request) {
	if !helpers.AlreadyLogin(r) {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	vars := mux.Vars(r)
	act := vars["action"]

	id, err := strconv.Atoi(r.FormValue("user_id"))
	if err != nil {
		log.Printf("Error converting string to int: %s", err)
	}

	switch act {
	case "make":
		{
			_, err = server.Core.DB.Exec("UPDATE profiles SET isadmin = $1 WHERE id = $2", true, id)
			if err != nil {
				log.Printf("Error giving %d's user admin rigths: %s", id, err)
			}
		}
	case "remove":
		{
			_, err = server.Core.DB.Exec("UPDATE profiles SET isadmin = $1 WHERE id = $2", false, id)
			if err != nil {
				log.Printf("Error removing %d's user admin rigths: %s", id, err)
			}
		}
	default:
		log.Printf("Invalid action")
	}
}

//GetItemHandler handle get requests
func GetItemHandler(w http.ResponseWriter, r *http.Request) {
	if !helpers.AlreadyLogin(r) {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	vars := mux.Vars(r)
	item := vars["item"]

	switch item {
	case "taskdata":
		{
			id, err := strconv.Atoi(r.FormValue("task_id"))
			if err != nil {
				log.Printf("Error converting task id(%d) to int: %s", id, err)
			}

			task, err := db.GetTaskbyID(server.Core.DB, id)
			if err != nil {
				log.Printf("Error getting task(%d) info: %s", id, err)
			}

			task.TaskExecutorName = strconv.Itoa(task.TaskExecutorID)

			formatTimeStart, err := time.Parse("02-01-2006", task.TaskDateStart)
			if err != nil {
				log.Printf("Error parsing date start for %s task for send to change page: %s", task.TaskName, err)
			}

			task.TaskDateStart = formatTimeStart.Format("2006-01-02")

			if task.TaskDateEnd != "" {
				formatTimeEnd, err := time.Parse("02-01-2006", task.TaskDateEnd)
				if err != nil {
					log.Printf("Error parsing date end for %s task for send to change page: %s", task.TaskName, err)
				}

				task.TaskDateEnd = formatTimeEnd.Format("2006-01-02")
			}

			taskData, err := json.MarshalIndent(task, "", " ")
			if err != nil {
				log.Printf("Error marshalling JSON for %s task: %s", task.TaskName, err)
			}

			w.Header().Set("Content-Type", "application/json")
			w.Write(taskData)
		}
	case "tmpldata":
		{
			id, err := strconv.Atoi(r.FormValue("tmpl_id"))
			if err != nil {
				log.Printf("Error converting template id(%d) to int: %s", id, err)
			}

			tmpl, err := db.GetTemplatebyID(server.Core.DB, id)
			if err != nil {
				log.Printf("Error getting template(%d) info: %s", id, err)
			}

			tmpl.TmplExecType = strconv.Itoa(tmpl.TmplExec)

			tmplData, err := json.MarshalIndent(tmpl, "", " ")
			if err != nil {
				log.Printf("Error marshalling JSON for %s template: %s", tmpl.TmplName, err)
			}

			w.Header().Set("Content-Type", "application/json")
			w.Write(tmplData)
		}
	case "taskdesc":
		{
			id := r.FormValue("id")
			path := "private/docs/" + id + ".md"

			bytes, err := ioutil.ReadFile(path)
			if err != nil {
				log.Printf("Error reading file with description for task with %s id: %s", id, err)
				w.Write([]byte("Ошибка при чтении файла с описанием: " + fmt.Sprintf("%s", err)))
			}

			w.Write(bytes)
		}
	case "wikiarticle":
		{
			id := r.FormValue("id")
			path := "private/wiki/" + id + ".md"

			bytes, err := ioutil.ReadFile(path)
			if err != nil {
				log.Printf("Error reading file with article for wiki with %s id: %s", id, err)
				w.Write([]byte("Ошибка при чтении файла со статьей: " + fmt.Sprintf("%s", err)))
			}

			w.Write(bytes)
		}
	case "wikilist":
		{
			wikies, err := db.GetAllWiki(server.Core.DB)
			if err != nil {
				log.Printf("Error getting wiki list: %s", err)
			}

			ans, err := json.MarshalIndent(wikies, "", "  ")
			if err != nil {
				log.Println("Error marshalling data to send: ", err)
			}

			w.Header().Set("Content-Type", "application/json")
			w.Write(ans)
		}
	}
}

//UserProfileHandler handle user's profile page
func UserProfileHandler(w http.ResponseWriter, r *http.Request) {
	if !helpers.AlreadyLogin(r) {
		http.Redirect(w, r, "/login", http.StatusMovedPermanently)
		return
	}

	viewData, err := db.GetDefaultViewData(server.Core.DB, r)
	if err != nil {
		log.Print("Error getting default viewData: ", err)
	}

	vars := mux.Vars(r)
	username := vars["user"]

	viewData.UserData.UserID, err = helpers.ConvertExecNameToID(username, "user")
	if err != nil {
		log.Print("Error getting user id: ", err)
	}

	viewData.UserData, err = db.GetUserbyID(server.Core.DB, viewData.UserData.UserID)
	if err != nil {
		log.Printf("Error getting info about %s: %s", username, err)
	}

	viewData.UserGroups, err = db.GetAllUserGroups(server.Core.DB, viewData.UserData.UserID)
	if err != nil {
		log.Printf("Error getting user's groups: %s", err)
	}

	viewData.UserExecTasks, err = db.GetAllTasksbyExecutor(server.Core.DB, viewData.UserData.UserID)
	if err != nil {
		log.Printf("Error getting assigned to user tasks: %s", err)
	}

	viewData.UserCreatorTasks, err = db.GetAllTasksbyCreator(server.Core.DB, viewData.UserData.UserID)
	if err != nil {
		log.Printf("Error getting created by user tasks: %s", err)
	}

	err = server.Core.Rnd.HTML(w, http.StatusOK, "profile", viewData)
	if err != nil {
		log.Print("Error parse template: ", err)
	}
}

//GroupsHandler handle page with all groups
func GroupsHandler(w http.ResponseWriter, r *http.Request) {
	if !helpers.AlreadyLogin(r) {
		http.Redirect(w, r, "/login", http.StatusMovedPermanently)
		return
	}

	viewData, err := db.GetDefaultViewData(server.Core.DB, r)
	if err != nil {
		log.Print("Error getting default viewData: ", err)
	}

	err = server.Core.Rnd.HTML(w, http.StatusOK, "groups", viewData)
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

	viewData, err := db.GetDefaultViewData(server.Core.DB, r)
	if err != nil {
		log.Print("Error getting default viewData: ", err)
	}

	vars := mux.Vars(r)
	groupIDstr := vars["id"]

	groupIDint, err := strconv.Atoi(groupIDstr)
	if err != nil {
		log.Printf("Error converting string to int on %s group page: %s", groupIDstr, err)
	}

	viewData.GroupData, err = db.GetGroupbyID(server.Core.DB, groupIDint)
	if err != nil {
		log.Printf("Error getting info about %d group: %s", groupIDint, err)
	}

	err = server.Core.Rnd.HTML(w, http.StatusOK, "group", viewData)
	if err != nil {
		log.Print("Error parse template: ", err)
	}
}

//TasksHandler handle page with tasks
func TasksHandler(w http.ResponseWriter, r *http.Request) {
	if !helpers.AlreadyLogin(r) {
		http.Redirect(w, r, "/login", http.StatusMovedPermanently)
		return
	}

	viewData, err := db.GetDefaultViewData(server.Core.DB, r)
	if err != nil {
		log.Print("Error getting default viewData: ", err)
	}

	err = server.Core.Rnd.HTML(w, http.StatusOK, "tasks", viewData)
	if err != nil {
		log.Print("Error parse template: ", err)
	}
}

//TaskPageHandler handle page of task
func TaskPageHandler(w http.ResponseWriter, r *http.Request) {
	if !helpers.AlreadyLogin(r) {
		http.Redirect(w, r, "/login", http.StatusMovedPermanently)
		return
	}

	viewData, err := db.GetDefaultViewData(server.Core.DB, r)
	if err != nil {
		log.Print("Error getting default viewData: ", err)
	}

	vars := mux.Vars(r)
	taskIDstr := vars["id"]

	taskIDint, err := strconv.Atoi(taskIDstr)
	if err != nil {
		log.Printf("Error converting string to int on %s task page: %s", taskIDstr, err)
	}

	viewData.TaskData, err = db.GetTaskbyID(server.Core.DB, taskIDint)
	if err != nil {
		log.Printf("Error getting task info from DB on %s task page: %s", taskIDstr, err)
	}

	err = server.Core.Rnd.HTML(w, http.StatusOK, "taskPage", viewData)
	if err != nil {
		log.Print("Error parse template: ", err)
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

	switch task.TaskExecutorType {
	case "user":
		{
			_, err = server.Core.DB.Exec("UPDATE profiles SET rating = (rating + $1) WHERE user_fio = $2", task.TaskRate, task.TaskExecutorFIO)
			if err != nil {
				log.Print(err)
			}
		}
	case "group":
		{
			group, err := db.GetGroupbyID(server.Core.DB, task.TaskExecutorID)
			if err != nil {
				log.Print(err)
			}

			rate := task.TaskRate / group.GroupMembersCount

			for _, user := range group.GroupMembers {
				_, err = server.Core.DB.Exec("UPDATE profiles SET rating = (rating + $1) WHERE user_fio = $2", rate, user.UserFIO)
				if err != nil {
					log.Print(err)
				}
			}
		}
	default:
		log.Printf("Error updating rate for group members: %s", errors.New("Wrong ExecutorType"))
	}

	_, err = server.Core.DB.Exec("UPDATE tasks SET stat = 'Закрыта' WHERE id = $1", task.TaskID)
	if err != nil {
		log.Print(err)
	}
}

//TaskOpenHandler handle page to reopen task
func TaskOpenHandler(w http.ResponseWriter, r *http.Request) {
	if !helpers.AlreadyLogin(r) {
		http.Redirect(w, r, "/login", http.StatusMovedPermanently)
		return
	}

	taskID, _ := strconv.Atoi(r.FormValue("id"))
	task, err := db.GetTaskbyID(server.Core.DB, taskID)
	if err != nil {
		log.Print(err)
	}

	_, err = server.Core.DB.Exec("UPDATE tasks SET (stat, rating) = ($1, $2) WHERE id = $3", "В процессе", 0, task.TaskID)
	if err != nil {
		log.Print("Can't update task: ", err)
	}
}

//WikiHandler handle page to wiki
func WikiHandler(w http.ResponseWriter, r *http.Request) {
	if !helpers.AlreadyLogin(r) {
		http.Redirect(w, r, "/login", http.StatusMovedPermanently)
		return
	}

	viewData, err := db.GetDefaultViewData(server.Core.DB, r)
	if err != nil {
		log.Print("Error getting viewData: ", err)
	}

	err = server.Core.Rnd.HTML(w, http.StatusOK, "wiki", viewData)
	if err != nil {
		log.Print("Error parse template: ", err)
	}
}

//WikiPageHandler handle wiki page
func WikiPageHandler(w http.ResponseWriter, r *http.Request) {
	if !helpers.AlreadyLogin(r) {
		http.Redirect(w, r, "/login", http.StatusMovedPermanently)
		return
	}

	viewData, err := db.GetDefaultViewData(server.Core.DB, r)
	if err != nil {
		log.Print("Error getting default viewData: ", err)
	}

	vars := mux.Vars(r)
	wikiIDstr := vars["id"]

	wikiIDint, err := strconv.Atoi(wikiIDstr)
	if err != nil {
		log.Printf("Error converting string to int on %s task page: %s", wikiIDstr, err)
	}

	viewData.WikiData, err = db.GetWikibyID(server.Core.DB, wikiIDint)
	if err != nil {
		log.Print("Error getting wiki data: ", err)
	}

	err = server.Core.Rnd.HTML(w, http.StatusOK, "wikiPage", viewData)
	if err != nil {
		log.Print("Error parse template: ", err)
	}
}
