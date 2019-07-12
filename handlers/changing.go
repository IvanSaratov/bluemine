package handlers

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

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
