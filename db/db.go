package db

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/IvanSaratov/bluemine/config"
	"github.com/IvanSaratov/bluemine/data"
	"github.com/IvanSaratov/bluemine/helpers"

	"github.com/go-ldap/ldap"
	"github.com/jmoiron/sqlx"
)

//RegisterUser adds user to DB
func RegisterUser(DB *sqlx.DB, l *ldap.Conn, login, userFIO string) (int64, error) {
	result, err := l.Search(ldap.NewSearchRequest(
		config.Conf.LdapBaseDN,
		ldap.ScopeWholeSubtree,
		ldap.NeverDerefAliases,
		0,
		0,
		false,
		"(&(sAMAccountName="+login+"))",
		[]string{"memberOf"},
		nil,
	))
	if err != nil {
		return 0, err
	}

	var listOfMembers []string
	for _, x := range result.Entries[0].GetAttributeValues("memberOf") {
		listOfMembers = append(listOfMembers, x[strings.Index(x, "CN=")+3:strings.Index(x, ",")])
	}

	isAdmin := false
	_, err = DB.Query("SELECT id FROM profiles")
	if err != nil {
		if err != sql.ErrNoRows {
			return 0, err
		}
		isAdmin = true
	}

	var userID int64
	err = DB.QueryRow("INSERT INTO profiles (username, user_fio, isAdmin) VALUES ($1, $2, $3) RETURNING id", login, userFIO, isAdmin).Scan(&userID)
	if err != nil {
		return 0, err
	}

	rows, err := DB.Query("SELECT id, group_name FROM groups")
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var groupID int
		var groupName string
		if err = rows.Scan(&groupID, &groupName); err != nil {
			return 0, err
		}
		for _, ldapGroupName := range listOfMembers {
			if groupName == ldapGroupName {
				_, err = DB.Exec("INSERT INTO groups_profiles (group_id, profile_id) VALUES ($1, $2)", groupID, userID)
				if err != nil {
					return 0, err
				}
				break
			}
		}
	}

	return userID, nil
}

//GetUserbyID gets user info from DB
func GetUserbyID(DB *sqlx.DB, id int) (data.User, error) {
	var (
		user data.User
		stmt = "SELECT * FROM profiles WHERE id = $1"
	)

	err := DB.QueryRow(stmt, id).Scan(&user.UserID, &user.UserName, &user.UserFIO, &user.UserisAdmin, &user.UserRate)
	if err != nil {
		return user, err
	}

	return user, nil
}

//GetAllUsers gets all users from DB
func GetAllUsers(DB *sqlx.DB) ([]data.User, error) {
	var (
		users []data.User
		stmt  = "SELECT id FROM profiles"
	)

	rows, err := DB.Query(stmt)
	if err != nil {
		return users, err
	}
	defer rows.Close()

	for rows.Next() {
		var user data.User

		err = rows.Scan(&user.UserID)
		if err != nil {
			return users, err
		}

		user, err = GetUserbyID(DB, user.UserID)
		if err != nil {
			return users, err
		}

		users = append(users, user)
	}
	if rows.Err() != nil {
		return users, err
	}

	return users, nil
}

//GetTaskbyID gets task info from DB
func GetTaskbyID(DB *sqlx.DB, ID int) (data.Task, error) {
	var (
		task data.Task
		stmt = "SELECT * FROM tasks WHERE id = $1"
	)

	err := DB.QueryRow(stmt, ID).Scan(&task.TaskID, &task.TaskName, &task.TaskCreatorID, &task.TaskExecutorID, &task.TaskExecutorType, &task.TaskStat, &task.TaskPriority, &task.TaskDateAdded, &task.TaskDateLastUpdate, &task.TaskDateStart, &task.TaskDateEnd, &task.TaskRate)
	if err != nil {
		return task, err
	}

	task.TaskCreatorName, err = helpers.ConvertIDToExecName(task.TaskCreatorID, "user")
	if err != nil {
		return task, err
	}

	task.TaskCreatorFIO, err = helpers.ConvertIDToExecFIO(task.TaskCreatorID)
	if err != nil {
		return task, err
	}

	task.TaskExecutorName, err = helpers.ConvertIDToExecName(task.TaskExecutorID, task.TaskExecutorType)
	if err != nil {
		return task, err
	}

	if task.TaskDateEnd != "" {
		timeEnd, err := time.Parse("02-01-2006", task.TaskDateEnd)
		if err != nil {
			log.Printf("Error parsing date end for %s task for calculate difference: %s", task.TaskName, err)
		}

		duration := timeEnd.Sub(time.Now())
		task.TaskDateDiff = duration.Hours()
	}

	if task.TaskExecutorType == "user" {
		task.TaskExecutorFIO, err = helpers.ConvertIDToExecFIO(task.TaskExecutorID)
		if err != nil {
			return task, err
		}
	}

	task.TaskChecklist, err = GetTaskCheckboxes(DB, task.TaskID)
	if err != nil {
		return task, err
	}

	return task, nil
}

//GetAllTasks gets all tasks from DB
func GetAllTasks(DB *sqlx.DB) ([]data.Task, error) {
	var (
		tasks []data.Task
		stmt  = "SELECT id FROM tasks"
	)

	rows, err := DB.Query(stmt)
	if err != nil {
		return tasks, err
	}
	defer rows.Close()

	for rows.Next() {
		var task data.Task

		err = rows.Scan(&task.TaskID)
		if err != nil {
			return tasks, err
		}

		task, err = GetTaskbyID(DB, task.TaskID)
		if err != nil {
			return tasks, err
		}

		tasks = append(tasks, task)
	}
	if rows.Err() != nil {
		return tasks, err
	}

	return tasks, nil
}

//GetTaskCheckboxes gets all task's checkboxes
func GetTaskCheckboxes(DB *sqlx.DB, ID int) ([]data.Checkbox, error) {
	var (
		checkboxes []data.Checkbox
		stmt       = "SELECT * from checkboxes WHERE task_id = $1"
	)

	rows, err := DB.Query(stmt, ID)
	if err != nil {
		return checkboxes, err
	}
	defer rows.Close()

	for rows.Next() {
		var checkbox data.Checkbox

		err = rows.Scan(&checkbox.CheckboxID, &checkbox.TaskID, &checkbox.Checked, &checkbox.CheckName)
		if err != nil {
			return checkboxes, err
		}

		checkboxes = append(checkboxes, checkbox)
	}
	if rows.Err() != nil {
		return checkboxes, err
	}

	return checkboxes, nil
}

//GetAllTasksbyExecutor gets all task with this executor
func GetAllTasksbyExecutor(DB *sqlx.DB, ID int) ([]data.Task, error) {
	var (
		tasks []data.Task
		stmt  = "SELECT id FROM tasks WHERE executor_id = $1"
	)

	rows, err := DB.Query(stmt, ID)
	if err != nil {
		return tasks, err
	}
	defer rows.Close()

	for rows.Next() {
		var task data.Task

		err = rows.Scan(&task.TaskID)
		if err != nil {
			return tasks, err
		}

		task, err = GetTaskbyID(DB, task.TaskID)
		if err != nil {
			return tasks, err
		}

		tasks = append(tasks, task)
	}
	if rows.Err() != nil {
		return tasks, err
	}

	userGroups, err := GetAllUserGroups(DB, ID)
	if err != nil {
		return tasks, err
	}

	for _, group := range userGroups {
		rows, err = DB.Query(stmt, group.GroupID)
		if err != nil {
			return tasks, err
		}
		defer rows.Close()

		for rows.Next() {
			var task data.Task

			err = rows.Scan(&task.TaskID)
			if err != nil {
				return tasks, err
			}

			task, err = GetTaskbyID(DB, task.TaskID)
			if err != nil {
				return tasks, err
			}

			tasks = append(tasks, task)
		}
		if rows.Err() != nil {
			return tasks, err
		}
	}

	return tasks, nil
}

//GetAllTasksbyCreator gets all task with this creator
func GetAllTasksbyCreator(DB *sqlx.DB, ID int) ([]data.Task, error) {
	var (
		tasks []data.Task
		stmt  = "SELECT id FROM tasks WHERE task_creator = $1"
	)

	rows, err := DB.Query(stmt, ID)
	if err != nil {
		return tasks, err
	}
	defer rows.Close()

	for rows.Next() {
		var task data.Task

		err = rows.Scan(&task.TaskID)
		if err != nil {
			return tasks, err
		}

		task, err = GetTaskbyID(DB, task.TaskID)
		if err != nil {
			return tasks, err
		}

		tasks = append(tasks, task)
	}
	if rows.Err() != nil {
		return tasks, err
	}

	return tasks, nil
}

//GetTemplatebyID gets task template info from DB
func GetTemplatebyID(DB *sqlx.DB, ID int) (data.TaskTmpl, error) {
	var (
		tmpl data.TaskTmpl
		stmt = "SELECT * FROM templates WHERE id = $1"
	)

	err := DB.QueryRow(stmt, ID).Scan(&tmpl.TmplID, &tmpl.TmplName, &tmpl.TmplExec, &tmpl.TmplExecType, &tmpl.TmplPriority, &tmpl.TmplRate)
	if err != nil {
		return tmpl, err
	}

	return tmpl, nil
}

//GetAllTemplates gets all task templates from DB
func GetAllTemplates(DB *sqlx.DB) ([]data.TaskTmpl, error) {
	var (
		tmpls []data.TaskTmpl
		stmt  = "SELECT id FROM templates"
	)

	rows, err := DB.Query(stmt)
	if err != nil {
		return tmpls, err
	}
	defer rows.Close()

	for rows.Next() {
		var tmpl data.TaskTmpl

		err = rows.Scan(&tmpl.TmplID)
		if err != nil {
			return tmpls, err
		}

		tmpl, err = GetTemplatebyID(DB, tmpl.TmplID)
		if err != nil {
			return tmpls, err
		}

		tmpls = append(tmpls, tmpl)
	}
	if rows.Err() != nil {
		return tmpls, err
	}

	return tmpls, nil
}

//GetGroupbyID gets group info from DB
func GetGroupbyID(DB *sqlx.DB, ID int) (data.Group, error) {
	var (
		group data.Group
		stmt  = "SELECT * FROM groups WHERE id = $1"
	)

	err := DB.QueryRow(stmt, ID).Scan(&group.GroupID, &group.GroupName)
	if err != nil {
		return group, err
	}

	group.GroupMembers, err = GetGroupUsers(DB, ID)
	if err != nil {
		return group, err
	}

	group.GroupMembersCount = len(group.GroupMembers)

	for _, member := range group.GroupMembers {
		group.GroupRate += member.UserRate
	}

	return group, nil
}

//GetAllGroups gets all users from DB
func GetAllGroups(DB *sqlx.DB) ([]data.Group, error) {
	var (
		groups []data.Group
		stmt   = "SELECT id FROM groups"
	)

	rows, err := DB.Query(stmt)
	if err != nil {
		return groups, err
	}
	defer rows.Close()

	for rows.Next() {
		var group data.Group

		err = rows.Scan(&group.GroupID)
		if err != nil {
			return groups, err
		}

		group, err = GetGroupbyID(DB, group.GroupID)
		if err != nil {
			return groups, err
		}

		groups = append(groups, group)
	}
	if rows.Err() != nil {
		return groups, err
	}

	return groups, nil
}

//GetGroupUsers gets all users of group from DB
func GetGroupUsers(DB *sqlx.DB, groupID int) ([]data.User, error) {
	var (
		users []data.User
		stmt  = "SELECT profile_id FROM groups_profiles WHERE group_id = $1"
	)

	rows, err := DB.Query(stmt, groupID)
	if err != nil {
		return users, err
	}
	defer rows.Close()

	for rows.Next() {
		var user data.User

		err = rows.Scan(&user.UserID)
		if err != nil {
			return users, err
		}

		user, err = GetUserbyID(DB, user.UserID)
		if err != nil {
			return users, err
		}

		users = append(users, user)
	}

	return users, nil
}

//GetAllUserGroups gets all groups of one users DB
func GetAllUserGroups(DB *sqlx.DB, ID int) ([]data.Group, error) {
	var groups []data.Group

	rows, err := DB.Query("SELECT group_id FROM groups_profiles WHERE profile_id = $1", ID)
	if err != nil {
		return groups, err
	}
	defer rows.Close()

	for rows.Next() {
		var group data.Group

		err = rows.Scan(&group.GroupID)
		if err != nil {
			return groups, err
		}

		group, err = GetGroupbyID(DB, group.GroupID)
		if err != nil {
			return groups, err
		}

		groups = append(groups, group)
	}

	return groups, nil
}

//GetWikibyID gets wiki by ID
func GetWikibyID(DB *sqlx.DB, ID int) (data.Wiki, error) {
	var (
		wiki data.Wiki
		stmt = "SELECT * FROM wiki WHERE id = $1"
	)

	err := DB.QueryRow(stmt, ID).Scan(&wiki.WikiID, &wiki.WikiAuthor.UserID, &wiki.WikiFatherID, &wiki.WikiName)
	if err != nil {
		return wiki, err
	}

	wiki.WikiIDStr = strconv.Itoa(wiki.WikiID)
	wiki.WikiFatherIDStr = strconv.Itoa(wiki.WikiFatherID)

	wiki.WikiAuthor, err = GetUserbyID(DB, wiki.WikiAuthor.UserID)
	if err != nil {
		return wiki, err
	}

	return wiki, nil
}

//GetAllWiki gets all wiki page
func GetAllWiki(DB *sqlx.DB) ([]data.Wiki, error) {
	var wikies []data.Wiki

	rows, err := DB.Query("SELECT id FROM wiki")
	if err != nil {
		return wikies, err
	}
	defer rows.Close()

	for rows.Next() {
		var wiki data.Wiki

		err = rows.Scan(&wiki.WikiID)
		if err != nil {
			return wikies, err
		}

		wiki, err = GetWikibyID(DB, wiki.WikiID)
		if err != nil {
			return wikies, err
		}

		wikies = append(wikies, wiki)
	}

	return wikies, nil
}

//GetDefaultViewData getting default view data
func GetDefaultViewData(DB *sqlx.DB, r *http.Request) (data.ViewData, error) {
	var viewData data.ViewData

	currentUser, err := helpers.GetCurrentUser(r)
	if err != nil {
		return viewData, err
	}

	wikies, err := GetAllWiki(DB)
	if err != nil {
		return viewData, err
	}

	tasks, err := GetAllTasks(DB)
	if err != nil {
		return viewData, err
	}

	users, err := GetAllUsers(DB)
	if err != nil {
		return viewData, err
	}

	groups, err := GetAllGroups(DB)
	if err != nil {
		return viewData, err
	}

	tmpls, err := GetAllTemplates(DB)
	if err != nil {
		return viewData, err
	}

	viewData = data.ViewData{
		CurrentUser: currentUser,
		Wikies:      wikies,
		Tasks:       tasks,
		Users:       users,
		Groups:      groups,
		Templates:   tmpls,
	}

	return viewData, nil
}

//DeleteRecord deleting user profile, template, wiki, task or groups
func DeleteRecord(DB *sqlx.DB, typeField string, ID int) error {
	var stmt string

	switch typeField {
	case "profile", "user":
		stmt = "DELETE FROM profiles WHERE id = $1"
	case "task":
		stmt = "DELETE FROM tasks WHERE id = $1"
	case "group":
		stmt = "DELETE FROM groups WHERE id = $1"
	case "template":
		stmt = "DELETE FROM templates WHERE id = $1"
	case "wiki":
		stmt = "DELETE FROM wiki WHERE id = $1"
	default:
		return errors.New("Wrong type")
	}

	_, err := DB.Exec(stmt, ID)
	if err != nil {
		return err
	}

	return nil
}
