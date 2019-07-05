package data

//User struct describe user
type User struct {
	UserID      int
	UserName    string
	UserFIO     string
	UserisAdmin bool
	UserRate    int
}

//Task struct describe tasks
type Task struct {
	TaskID           int
	TaskName         string
	TaskCreator      User
	TaskExecutor     User
	TaskExecutorType string
	TaskStat         string
	TaskDateStart    string
	TaskDateEnd      string
	TaskRate         int
}

//Group struct describe group
type Group struct {
	GroupID      int
	GroupName    string
	GroupMembers []User
}

//Wiki struct describe wiki article
type Wiki struct {
	WikiID       int
	WikiAuthor   User
	WikiFatherID int
	WikiName     string
}

//ViewData struct to show data on page
type ViewData struct {
	CurrentUser User
	UserData    User
	GroupData   Group
	TaskData    Task
	Users       []User
	Groups      []Group
	Tasks       []Task
}
