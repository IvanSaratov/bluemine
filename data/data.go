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
	TaskID        int
	TaskName      string
	TaskDescPath  string
	TaskExecutorType string
	TaskExecutor  string
	TaskStat      string
	TaskDateStart string
	TaskDateEnd   string
	TaskRate      int
}

//ViewData struct to show data on page
type ViewData struct {
	CurrentUser User
	UserData    User
	TaskData    Task
	Tasks       []Task
}
