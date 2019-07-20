package main

import (
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/IvanSaratov/bluemine/data"
	"github.com/IvanSaratov/bluemine/db"

	"github.com/IvanSaratov/bluemine/config"
	"github.com/IvanSaratov/bluemine/handlers"
	"github.com/IvanSaratov/bluemine/server"

	_ "github.com/cockroachdb/cockroach-go/crdb"
	"github.com/gorilla/mux"
)

func init() {
	var configPath = "conf.toml"

	log.Println("Starting...")

	err := config.ParseConfig(configPath)
	if err != nil {
		log.Fatal("Error parsing config: ", err)
	}
	log.Println("Config parsed!")

	server.Init()
}

func main() {
	defer server.Core.DB.Close()

	router := mux.NewRouter()

	router.PathPrefix("/private/").HandlerFunc(handlers.PrivateHandler)
	router.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("./public/"))))
	router.HandleFunc("/login", handlers.LoginHandler)
	router.HandleFunc("/logout", handlers.LogoutHandler)
	router.HandleFunc("/profile/{user}", handlers.UserProfileHandler)
	router.HandleFunc("/makeadmin", handlers.MakeAdminHandler).Methods("POST")
	router.HandleFunc("/removeadmin", handlers.RemoveAdminHandler).Methods("POST")
	router.HandleFunc("/group/show/{id}", handlers.GroupHandler)
	router.HandleFunc("/group/new", handlers.AddGroupHandler).Methods("POST")
	router.HandleFunc("/groups", handlers.GroupsHandler)
	router.HandleFunc("/group/change", handlers.GroupChangeHandler)
	router.HandleFunc("/tasks", handlers.TasksHandler)
	router.HandleFunc("/tasks/show/{id}", handlers.TaskPageHandler)
	router.HandleFunc("/gettaskdata", handlers.GetTaskData).Methods("GET")
	router.HandleFunc("/getdesc", handlers.GetPrivateDesc).Methods("GET")
	router.HandleFunc("/getwikilist", handlers.GetWikiTree)
	router.HandleFunc("/tasks/new", handlers.AddTaskHandler).Methods("POST")
	router.HandleFunc("/tasks/change", handlers.ChangeTaskHandler).Methods("POST")
	router.HandleFunc("/tmpl/new", handlers.AddTmplHandler).Methods("POST")
	router.HandleFunc("/gettmpldata", handlers.GetTmplData).Methods("GET")
	router.HandleFunc("/tasks/close", handlers.TaskCloseHandler)
	router.HandleFunc("/tasks/open", handlers.TaskOpenHandler).Methods("POST")
	router.HandleFunc("/wiki", handlers.WikiHandler)
	router.HandleFunc("/wiki/show/{id}", handlers.WikiPageHandler)
	router.HandleFunc("/wiki/new", handlers.AddWikiHandler)
	router.HandleFunc("/", handlers.RootHandler)

	go taskAutoCloser()

	log.Printf("Server listening on %s port", config.Conf.ListenPort)
	log.Fatal(http.ListenAndServe(config.Conf.ListenPort, router))
}

//taskAutoCloser close task automatically when the term expires
func taskAutoCloser() {
	for range time.Tick(time.Minute * 5) {
		tasks, err := db.GetAllTasks(server.Core.DB)
		if err != nil {
			log.Printf("Error getting task list to auto close: %s", err)
		}

		for _, task := range tasks {
			go func(task data.Task) {
				if task.TaskDateEnd != "" {
					if task.TaskDateDiff < float64(-3*time.Hour) {
						switch task.TaskExecutorType {
						case "user":
							{
								_, err = server.Core.DB.Exec("UPDATE profiles SET rating = (rating - $1) WHERE user_fio = $2", task.TaskRate, task.TaskExecutorFIO)
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
									_, err = server.Core.DB.Exec("UPDATE profiles SET rating = (rating - $1) WHERE user_fio = $2", rate, user.UserFIO)
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
				}
			}(task)
		}
	}
}
