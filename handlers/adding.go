package handlers

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/IvanSaratov/bluemine/data"
	"github.com/IvanSaratov/bluemine/db"
	"github.com/IvanSaratov/bluemine/helpers"
	"github.com/IvanSaratov/bluemine/server"
)

//AddTaskHandler handle task adding
func AddTaskHandler(w http.ResponseWriter, r *http.Request) {
	if !helpers.AlreadyLogin(r) {
		http.Redirect(w, r, "/login", http.StatusMovedPermanently)
		return
	}

	currentUser, err := helpers.GetCurrentUser(r)
	if err != nil {
		log.Printf("Error getting current user: %s", err)
	}

	if r.Method == "POST" {
		var (
			task        data.Task
			description string
			err         error
		)

		task.TaskName = r.FormValue("task_name")

		task.TaskCreatorID = currentUser.UserID

		task.TaskExecutorID, err = strconv.Atoi(r.FormValue("task_exec"))
		if err != nil {
			log.Printf("Error converting executor's ID to int: %s", err)
		}

		task.TaskExecutorType = r.FormValue("task_exec_type")

		task.TaskStat = r.FormValue("task_stat")

		task.TaskPriority = r.FormValue("task_priority")

		task.TaskDateAdded = time.Now().Format("2006-01-02 15:04:05")

		task.TaskDateLastUpdate = time.Now().Format("2006-01-02 15:04:05")

		task.TaskDateStart = r.FormValue("task_start")

		task.TaskDateEnd = r.FormValue("task_end")

		task.TaskRate, err = strconv.Atoi(r.FormValue("task_rate"))
		if err != nil {
			log.Printf("Error converting rating from string to int: %s", err)
		}

		description = r.FormValue("task_desc")

		err = server.Core.DB.QueryRow("INSERT INTO tasks (task_name, task_creator, executor_id, executor_type, stat, priority, date_added, date_last_update, date_start, date_end, rating) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) RETURNING id", task.TaskName, task.TaskCreatorID, task.TaskExecutorID, task.TaskExecutorType, task.TaskStat, task.TaskPriority, task.TaskDateAdded, task.TaskDateLastUpdate, task.TaskDateStart, task.TaskDateEnd, task.TaskRate).Scan(&task.TaskID)
		if err != nil {
			log.Print(err)
		}

		f, err := os.OpenFile("private/docs/"+strconv.Itoa(task.TaskID)+".txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Print(err)
		}

		_, err = f.WriteString(description)
		if err != nil {
			log.Print(err)
		}
	}
}

//AddTmplHandler handle template adding
func AddTmplHandler(w http.ResponseWriter, r *http.Request) {
	if !helpers.AlreadyLogin(r) {
		http.Redirect(w, r, "/login", http.StatusMovedPermanently)
		return
	}

	if r.Method == "POST" {
		var (
			tmpl data.TaskTmpl
			err  error
		)

		tmpl.TmplName = r.FormValue("tmpl_name")

		tmpl.TmplStat = r.FormValue("tmpl_stat")

		tmpl.TmplPriority = r.FormValue("tmpl_priority")

		tmpl.TmplRate, err = strconv.Atoi(r.FormValue("tmpl_rate"))
		if err != nil {
			log.Printf("Error converting rating from string to int: %s", err)
		}

		_, err = server.Core.DB.Exec("INSERT INTO task_template (tmpl_name, stat, priority, rating) VALUES ($1, $2, $3, $4) RETURNING id", tmpl.TmplName, tmpl.TmplStat, tmpl.TmplPriority, tmpl.TmplRate)
		if err != nil {
			log.Print(err)
		}
	}
}

//AddGroupHandler handle group create page
func AddGroupHandler(w http.ResponseWriter, r *http.Request) {
	if !helpers.AlreadyLogin(r) {
		http.Redirect(w, r, "/login", http.StatusMovedPermanently)
		return
	}

	currentUser, err := helpers.GetCurrentUser(r)
	if err != nil {
		log.Printf("Error getting current user: %s", err)
	}

	users, err := db.GetAllUsers(server.Core.DB)
	if err != nil {
		log.Printf("Error getting users list: %s", err)
	}

	if r.Method == "GET" {
		viewData := data.ViewData{
			CurrentUser: currentUser,
			Users:       users,
		}

		err := server.Core.Templates["addGroup"].ExecuteTemplate(w, "base", viewData)
		if err != nil {
			log.Print(err)
		}
	} else if r.Method == "POST" {
		var group data.Group

		group.GroupName = r.FormValue("input_group")

		for _, userString := range strings.Split(r.FormValue("user_list"), ",") {
			var user data.User
			user.UserID, err = strconv.Atoi(userString[strings.Index(userString, "user")+4:])
			if err != nil {
				log.Printf("Can't get userID: %s", err)
			}

			group.GroupMembers = append(group.GroupMembers, user)
		}

		err := server.Core.DB.QueryRow("INSERT INTO groups (group_name) VALUES ($1) RETURNING id", group.GroupName).Scan(&group.GroupID)
		if err != nil {
			log.Printf("Can't create group: %s", err)
		}

		for _, user := range group.GroupMembers {
			_, err := server.Core.DB.Exec("INSERT INTO groups_profiles (group_id, profile_id) VALUES ($1, $2)", group.GroupID, user.UserID)
			if err != nil {
				log.Printf("Can't create group-profile link: %s", err)
			}
		}
	}
}

//AddWikiHandler handle wiki adding page
/*func AddWikiHandler(w http.ResponseWriter, r *http.Request) {
	if !helpers.AlreadyLogin(r) {
		http.Redirect(w, r, "/login", 301)
		return
	}

	currentUser, err := helpers.GetCurrentUser(r)
	if err != nil {
		log.Printf("Error getting current user: %s", err)
	}

	if r.Method == "GET" {
		viewData := data.ViewData{
			CurrentUser: currentUser,
		}

		tmpl, _ := template.ParseFiles("public/html/addtask.html")
		tmpl.Execute(w, viewData)
	} else if r.Method == "POST" {
		var (
			wiki       data.Wiki
			authorName string
			err        error
		)

		wiki.WikiName = r.FormValue("wiki_name")

		session, _ := server.Core.Store.Get(r, "bluemine_session")
		authorName = fmt.Sprint(session.Values["userName"])
		wiki.WikiAuthorID, err = helpers.ConvertExecToID(authorName, "user")
		if err != nil {
			log.Print(err)
		}

		_, err = server.Core.DB.Exec("INSERT INTO wiki (wiki_name, author_id) VALUES ($1, $2)", wiki.WikiName, wiki.WikiAuthorID)
		if err != nil {
			log.Print(err)
		}
	}
}*/
