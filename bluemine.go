package main

import (
	"errors"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/IvanSaratov/bluemine/data"
	"github.com/IvanSaratov/bluemine/db"

	"github.com/IvanSaratov/bluemine/config"
	"github.com/IvanSaratov/bluemine/handlers"
	"github.com/IvanSaratov/bluemine/server"

	_ "github.com/cockroachdb/cockroach-go/crdb"
	"github.com/gorilla/mux"
)

var logFile *os.File

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

	router.HandleFunc("/admin/{action}", handlers.AdminActHandler).Methods("POST")

	router.HandleFunc("/group/show/{id}", handlers.GroupHandler)
	router.HandleFunc("/group/new", handlers.AddGroupHandler).Methods("POST")
	router.HandleFunc("/groups", handlers.GroupsHandler)
	router.HandleFunc("/group/change", handlers.GroupChangeHandler)
	router.HandleFunc("/tasks", handlers.TasksHandler)
	router.HandleFunc("/tasks/show/{id}", handlers.TaskPageHandler)
	router.HandleFunc("/tasks/new", handlers.AddTaskHandler).Methods("POST")
	router.HandleFunc("/tasks/change", handlers.ChangeTaskHandler).Methods("POST")
	router.HandleFunc("/tmpl/new", handlers.AddTmplHandler).Methods("POST")
	router.HandleFunc("/tasks/close", handlers.TaskCloseHandler)
	router.HandleFunc("/tasks/open", handlers.TaskOpenHandler).Methods("POST")
	router.HandleFunc("/wiki", handlers.WikiHandler)
	router.HandleFunc("/wiki/show/{id}", handlers.WikiPageHandler)
	router.HandleFunc("/wiki/new", handlers.AddWikiHandler)

	router.HandleFunc("/get/{item}", handlers.GetItemHandler).Methods("GET")

	router.HandleFunc("/", handlers.RootHandler)

	log.Printf("Server must listen on %s port", config.Conf.ListenPort)

	go logRotator()
	go taskAutoCloser()

	log.Fatal(http.ListenAndServe(config.Conf.ListenPort, router))
}

//logRotate creates log files
func logRotate() error {
	var logFilePath = "logs/" + time.Now().Format("2006-01-02") + ".log"
	if _, err := os.Stat(logFilePath); err != nil {
		if os.IsNotExist(err) {
			logFile, err := os.OpenFile("logs/"+time.Now().Format("2006-01-02")+".log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
			if err != nil {
				return err
			}

			log.SetOutput(logFile)
			log.Printf("Created log file %s", time.Now().Format("2006-01-02"))
		}
	}

	return nil
}

//logRotator rotate log files automatically
func logRotator() {
	err := logRotate()
	if err != nil {
		log.Printf("Error rotating %s log file: %s", time.Now().Format("2006-01-02"), err)
	}

	for range time.Tick(time.Hour * 24) {
		logFile.Close()
		err = logRotate()
		if err != nil {
			log.Printf("Error rotating %s log file: %s", time.Now().Format("2006-01-02"), err)
		}
	}
}

//taskAutoCloser close task automatically when the term expires
func taskAutoCloser() {
	var wg sync.WaitGroup

	for range time.Tick(time.Minute * 5) {
		tasks, err := db.GetAllTasks(server.Core.DB)
		if err != nil {
			log.Printf("Error getting task list to auto close: %s", err)
		}

		for _, task := range tasks {
			wg.Add(1)

			go func(task data.Task) {
				defer wg.Done()

				if task.TaskStat != "Закрыта" {
					if task.TaskDateEnd != "" {
						if task.TaskDateDiff < float64(-18) {
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
				}
			}(task)
		}
		wg.Wait()
	}
}
