package data

//User struct describe user
type User struct {
	UserID         int
	UserName       string
	UserFIO        string
	UserDepartment string
	UserGroup      string
	UserisAdmin    bool
	UserRate       int
}

//Task struct describe tasks
type Task struct {
	TaskID           int
	TaskName         string
	TaskExecutorType string
	TaskExecutor     string
	TaskStat         string
	TaskDateStart    string
	TaskDateEnd      string
	TaskRate         int
}

//Group struct describe group
type Group struct {
	GroupID   int
	GroupName string
}

//Wiki struct describe wiki article
type Wiki struct {
	WikiID       int
	WikiAuthorID int
	WikiFatherID int
	WikiName     string
}

//ViewData struct to show data on page
type ViewData struct {
	CurrentUser User
	UserData    User
	TaskData    Task
	Users       []User
	Tasks       []Task
}
