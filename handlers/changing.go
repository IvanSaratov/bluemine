package handlers

import (
	"log"
	"net/http"
	"encoding/json"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/IvanSaratov/bluemine/db"
	"github.com/IvanSaratov/bluemine/data"
	"github.com/IvanSaratov/bluemine/helpers"
	"github.com/IvanSaratov/bluemine/server"
)

//ChangeTaskHandler changes task
func ChangeTaskHandler(w http.ResponseWriter, r *http.Request) {
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

		task.TaskDateLastUpdate = time.Now().Format("2006-01-02 15:04:05")

		task.TaskDateStart = r.FormValue("task_start")

		task.TaskDateEnd = r.FormValue("task_end")

		task.TaskRate, err = strconv.Atoi(r.FormValue("task_rate"))
		if err != nil {
			log.Printf("Error converting rating from string to int: %s", err)
		}

		description = r.FormValue("task_desc")

		err = server.Core.DB.QueryRow("UPDATE tasks SET (task_name, task_creator, executor_id, executor_type, stat, priority, date_last_update, date_start, date_end, rating) = ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) RETURNING id", task.TaskName, task.TaskCreatorID, task.TaskExecutorID, task.TaskExecutorType, task.TaskStat, task.TaskPriority, task.TaskDateLastUpdate, task.TaskDateStart, task.TaskDateEnd, task.TaskRate).Scan(&task.TaskID)
		if err != nil {
			log.Print(err)
		}

		f, err := os.OpenFile("private/docs/"+strconv.Itoa(task.TaskID)+".txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Print(err)
		}

		err = f.Truncate(0)
		if err != nil {
			log.Println(err)
		}

		_, err = f.Seek(0, 0)
		if err != nil {
			log.Println(err)
		}

		_, err = f.WriteString(description)
		if err != nil {
			log.Print(err)
		}
	}
}

//GroupChangeHandler handler page to change group settings
func GroupChangeHandler(w http.ResponseWriter, r *http.Request) {
	if !helpers.AlreadyLogin(r) {
		http.Redirect(w, r, "/login", http.StatusMovedPermanently)
		return
	}

	if r.Method == "GET"{
	id, err := strconv.Atoi(r.FormValue("id"))
	if err != nil {
		log.Print("Error convert string to id: ", err)
	}

	group, err := db.GetGroupbyID(server.Core.DB, id)
	if err != nil {
		log.Print("Error getting group by id: ", err)
	}

	for x, user := range group.GroupMembers {
		userIDstring := strconv.Itoa(user.UserID)
		group.GroupMembers[x].UserName = userIDstring
	}

	groupData, err := json.MarshalIndent(group, "", " ")
	if err != nil {
		log.Printf("Error marshalling JSON for %s group: %s", group.GroupName, err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(groupData)
	} else if r.Method == "POST" {
		var (
			group data.Group
			err   error
		)

		if err := r.ParseForm(); err != nil {
			log.Printf("Error parsing form: %s", err)
		}

		for key, value := range r.PostForm {
			if key == "group_id" {
				group.GroupID, err = strconv.Atoi(value[0])
				if err != nil {
					log.Print("Error: ", err)
				}
			} else if key == "group_name" {
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

		_, err = server.Core.DB.Exec("UPDATE groups SET group_name = $1 WHERE id = $2", group.GroupName, group.GroupID)
		if err != nil {
			log.Print("Error set new group name: ", err)
		}

		_, err = server.Core.DB.Exec("DELETE FROM groups_profiles WHERE group_id = $1", group.GroupID)
		if err != nil {
			log.Print("Error deleting group users: ", err)
		}

		for _, user := range group.GroupMembers {
			_, err := server.Core.DB.Exec("INSERT INTO groups_profiles (group_id, profile_id) VALUES ($1, $2)", group.GroupID, user.UserID)
			if err != nil {
				log.Printf("Error writing groups and users to groups_profiles table: %s", err)
			}
		}
	}
}