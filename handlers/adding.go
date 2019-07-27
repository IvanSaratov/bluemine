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

	var (
		task        data.Task
		description string
		checklist   []string
	)

	currentUser, err := helpers.GetCurrentUser(r)
	if err != nil {
		log.Printf("Error getting current user: %s", err)
	}

	task.TaskName = r.FormValue("task_name")

	task.TaskCreatorID = currentUser.UserID

	task.TaskExecutorID, err = strconv.Atoi(r.FormValue("task_exec"))
	if err != nil {
		log.Printf("Error converting executor's ID to int: %s", err)
	}

	task.TaskExecutorType = r.FormValue("task_exec_type")

	task.TaskStat = r.FormValue("task_stat")

	task.TaskPriority = r.FormValue("task_priority")

	task.TaskDateAdded = time.Now().Format("02-01-2006 15:04:05")

	task.TaskDateLastUpdate = time.Now().Format("02-01-2006 15:04:05")

	timeStart, err := time.Parse("2006-01-02", r.FormValue("task_start"))
	if err != nil {
		log.Printf("Error parsing date start for %s task: %s", task.TaskName, err)
	}

	task.TaskDateStart = timeStart.Format("02-01-2006")

	if r.FormValue("task_end") == "" {
		task.TaskDateEnd = ""
	} else {
		timeEnd, err := time.Parse("2006-01-02", r.FormValue("task_end"))
		if err != nil {
			log.Printf("Error parsing date end for %s task: %s", task.TaskName, err)
		}

		task.TaskDateEnd = timeEnd.Format("02-01-2006")
	}

	task.TaskRate, err = strconv.Atoi(r.FormValue("task_rate"))
	if err != nil {
		log.Printf("Error converting rating from string to int: %s", err)
	}

	description = r.FormValue("task_desc")

	checklist = strings.Split(r.FormValue("task_checklist"), "&")
	for _, checkboxStr := range checklist {
		var (
			checkbox data.Checkbox
			checkmap = strings.Split(checkboxStr, "=")
		)

		checkbox.CheckName = checkmap[0]
		checkbox.Checked, err = strconv.ParseBool(checkmap[1])
		if err != nil {
			log.Printf("Error setting %s checkbox check status for %s task: %s", checkbox.CheckName, task.TaskName, err)
		}

		_, err = server.Core.DB.Exec("INSERT INTO checkboxes (task_id, checked, desk) VALUES ($1, $2, $3)", task.TaskID, checkbox.Checked, checkbox.CheckName)
		if err != nil {
			log.Printf("Error inserting %s checkbox into DB for %s task: %s", checkbox.CheckName, task.TaskName, err)
		}
	}

	err = server.Core.DB.QueryRow("INSERT INTO tasks (task_name, task_creator, executor_id, executor_type, stat, priority, date_added, date_last_update, date_start, date_end, rating) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) RETURNING id", task.TaskName, task.TaskCreatorID, task.TaskExecutorID, task.TaskExecutorType, task.TaskStat, task.TaskPriority, task.TaskDateAdded, task.TaskDateLastUpdate, task.TaskDateStart, task.TaskDateEnd, task.TaskRate).Scan(&task.TaskID)
	if err != nil {
		log.Print(err)
	}

	f, err := os.OpenFile("private/docs/"+strconv.Itoa(task.TaskID)+".md", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Print(err)
	}

	_, err = f.WriteString(description)
	if err != nil {
		log.Print(err)
	}
}

//AddTmplHandler handle template adding
func AddTmplHandler(w http.ResponseWriter, r *http.Request) {
	if !helpers.AlreadyLogin(r) {
		http.Redirect(w, r, "/login", http.StatusMovedPermanently)
		return
	}

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

	_, err = server.Core.DB.Exec("INSERT INTO templates (tmpl_name, stat, priority, rating) VALUES ($1, $2, $3, $4) RETURNING id", tmpl.TmplName, tmpl.TmplStat, tmpl.TmplPriority, tmpl.TmplRate)
	if err != nil {
		log.Print(err)
	}
}

//AddGroupHandler handle group create page
func AddGroupHandler(w http.ResponseWriter, r *http.Request) {
	if !helpers.AlreadyLogin(r) {
		http.Redirect(w, r, "/login", http.StatusMovedPermanently)
		return
	}

	var (
		group data.Group
		err   error
	)

	if err := r.ParseForm(); err != nil {
		log.Printf("Error parsing form: %s", err)
	}

	for key, value := range r.PostForm {
		if key == "group_name" {
			group.GroupName = value[0]
		} else if key == "user_list" {
			users := strings.Split(value[0], "&")
			for _, userID := range users {
				var user data.User

				user.UserID, err = strconv.Atoi(userID[5:])
				if err != nil {
					log.Printf("Can't get userID: %s", err)
				}

				group.GroupMembers = append(group.GroupMembers, user)
			}
		}
	}

	err = server.Core.DB.QueryRow("INSERT INTO groups (group_name) VALUES ($1) RETURNING id", group.GroupName).Scan(&group.GroupID)
	if err != nil {
		log.Printf("Error writing group data to DB: %s", err)
	}

	for _, user := range group.GroupMembers {
		_, err := server.Core.DB.Exec("INSERT INTO groups_profiles (group_id, profile_id) VALUES ($1, $2)", group.GroupID, user.UserID)
		if err != nil {
			log.Printf("Error writing groups and users to groups_profiles table: %s", err)
		}
	}
}

//AddWikiHandler handle wiki adding page
func AddWikiHandler(w http.ResponseWriter, r *http.Request) {
	if !helpers.AlreadyLogin(r) {
		http.Redirect(w, r, "/login", http.StatusMovedPermanently)
		return
	}

	if r.Method == "GET" {
		viewData, err := db.GetDefaultViewData(server.Core.DB, r)
		if err != nil {
			log.Print("Error getting default viewData: ", err)
		}

		err = server.Core.Templates["addWiki"].ExecuteTemplate(w, "index", viewData)
		if err != nil {
			log.Print(err)
		}

	} else if r.Method == "POST" {
		var (
			wiki    data.Wiki
			err     error
			article string
		)

		wiki.WikiName = r.FormValue("wiki_name")
		article = r.FormValue("article")
		wiki.WikiAuthor, err = helpers.GetCurrentUser(r)
		if err != nil {
			log.Print("Error getting current user: ", err)
		}

		wiki.WikiFatherID, err = strconv.Atoi(r.FormValue("father_id"))
		if err != nil {
			log.Print("Error getting father_id: ", err)
		}

		err = server.Core.DB.QueryRow("INSERT INTO wiki (title, author_id, father_id) VALUES ($1, $2, $3) RETURNING id", wiki.WikiName, wiki.WikiAuthor.UserID, wiki.WikiFatherID).Scan(&wiki.WikiID)
		if err != nil {
			log.Print("Error create wiki article: ", err)
		}

		f, err := os.OpenFile("private/wiki/"+strconv.Itoa(wiki.WikiID)+".md", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Print("Error create wiki.md: ", err)
		}

		_, err = f.WriteString(article)
		if err != nil {
			log.Print("Error to write in md: ", err)
		}
	}
}
