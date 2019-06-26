package data

//User struct describe user
type User struct {
	UserID         int
	UserName       string
	UserFIO        string
	UserDepartment int
	UserGroup      int
	UserisAdmin    bool
	UserRate       int
}

//Task struct describe tasks
type Task struct {
	TaskID        int
	TaskName      string
	TaskDesc      string
	TaskStat      int
	TaskDateStart string
	TaskDateEnd   string
	TaskRate      int
}

//ViewData struct to show data on page
type ViewData struct {
	UserData User
	TaskData Task
}
