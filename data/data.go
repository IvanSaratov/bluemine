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
	TaskID   int
	TaskName string

	TaskCreatorID   int
	TaskCreatorName string
	TaskCreatorFIO  string

	TaskExecutorID   int
	TaskExecutorName string
	TaskExecutorFIO  string
	TaskExecutorType string

	TaskStat           string
	TaskPriority       string
	TaskDateAdded      string
	TaskDateLastUpdate string
	TaskDateStart      string
	TaskDateEnd        string
	TaskRate           int
}

//TaskTmpl struct describe task template
type TaskTmpl struct {
	TmplID       int
	TmplName     string
	TmplStat     string
	TmplPriority string
	TmplRate     int
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
	CurrentUser      User
	UserData         User
	GroupData        Group
	TaskData         Task
	Users            []User
	Groups           []Group
	UserGroups       []Group
	Tasks            []Task
	UserExecTasks    []Task
	UserCreatorTasks []Task
	Templates        []TaskTmpl
}
