package handlers

import (
	"html/template"
	"net/http"
	"log"
	"strconv"

	"github.com/IvanSaratov/bluemine/data"
	"github.com/IvanSaratov/bluemine/helpers"
	"github.com/IvanSaratov/bluemine/server"
)

//AddTaskHandler handle task adding page
func AddTaskHandler(w http.ResponseWriter, r *http.Request) {
	if !helpers.AlreadyLogin(r) {
		http.Redirect(w, r, "/login", 301)
		return
	}

	if r.Method == "GET" {
		data := data.ViewData{
			UserData: data.User{
				UserName:       "test",
				UserFIO:        "test_testovich",
				UserDepartment: "Otdel_Debilov",
			},
		}
	
		tmpl, _ := template.ParseFiles("public/html/addtask.html")
		tmpl.Execute(w, data)
	} else if r.Method == "POST" {
		var (
			task data.Task
			executorID int
			description string
			err error
		)
		
		task.TaskName = r.FormValue("task_name")
		task.TaskStat = r.FormValue("task_stat")
		task.TaskDateStart = r.FormValue("task_start")
		task.TaskDateEnd = r.FormValue("task_end")
		task.TaskExecutorType = r.FormValue("executor_type")
		task.TaskExecutor = r.FormValue("executor_name")
		task.TaskRate, _ = strconv.Atoi(r.FormValue("task_rate"))
		description = r.FormValue("task_desc")

		executorID, err = helpers.ConvertExecToID(task.TaskExecutor, task.TaskExecutorType)
		if err != nil {
			log.Print(err)
		}

		err = server.Core.DB.QueryRow("INSERT INTO tasks (task_name, stat, date_start, date_end, rating, executor_type, executor_id) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id", task.TaskName, task.TaskStat, task.TaskDateStart, task.TaskDateEnd, task.TaskRate, task.TaskExecutorType, executorID).Scan(&task.TaskID)
		if err != nil {
			log.Print(err)
		}
	}
}
