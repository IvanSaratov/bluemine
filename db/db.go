package db

import (
	"log"
	"strings"
	"time"

	"github.com/IvanSaratov/bluemine/config"
	"github.com/IvanSaratov/bluemine/data"
	"github.com/IvanSaratov/bluemine/helpers"

	"github.com/go-ldap/ldap"
	"github.com/jmoiron/sqlx"
)

//RegisterUser adds user to DB
func RegisterUser(DB *sqlx.DB, l *ldap.Conn, login, userFIO string) error {
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
		return err
	}

	var listOfMembers []string
	for _, x := range result.Entries[0].GetAttributeValues("memberOf") {
		listOfMembers = append(listOfMembers, x[strings.Index(x, "CN=")+3:strings.Index(x, ",")])
	}

	var userID int64
	err = DB.QueryRow("INSERT INTO profiles (username, user_fio) VALUES ($1, $2) RETURNING id", login, userFIO).Scan(&userID)
	if err != nil {
		return err
	}

	rows, err := DB.Query("SELECT id, group_name FROM groups")
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var groupID int
		var groupName string
		if err = rows.Scan(&groupID, &groupName); err != nil {
			return err
		}
		for _, ldapGroupName := range listOfMembers {
			if groupName == ldapGroupName {
				_, err = DB.Exec("INSERT INTO groups_profiles (group_id, profile_id) VALUES ($1, $2)", groupID, userID)
				if err != nil {
					return err
				}
				break
			}
		}
	}

	return nil
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

//GetAllTasksbyExecutor gets all task with this executor
func GetAllTasksbyExecutor(DB *sqlx.DB, ID int) ([]data.Task, error) {
	var (
		tasks []data.Task
		stmt  = "SELECT id FROM tasks WHERE executor_type = $1 AND executor_id = $2"
	)

	rows, err := DB.Query(stmt, "user", ID)
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
		stmt = "SELECT * FROM task_template WHERE id = $1"
	)

	err := DB.QueryRow(stmt, ID).Scan(&tmpl.TmplID, &tmpl.TmplName, &tmpl.TmplStat, &tmpl.TmplPriority, &tmpl.TmplRate)
	if err != nil {
		return tmpl, err
	}

	return tmpl, nil
}

//GetAllTemplates gets all task templates from DB
func GetAllTemplates(DB *sqlx.DB) ([]data.TaskTmpl, error) {
	var (
		tmpls []data.TaskTmpl
		stmt  = "SELECT id FROM task_template"
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
