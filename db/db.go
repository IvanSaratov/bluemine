package db

import (
	"database/sql"
	"strings"

	"github.com/IvanSaratov/bluemine/helpers"
	"github.com/IvanSaratov/bluemine/data"
	"github.com/IvanSaratov/bluemine/config"

	"github.com/go-ldap/ldap"
)

//RegisterUser adds user to DB
func RegisterUser(DB *sql.DB, l *ldap.Conn, login, userFIO string) error {
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

//GetUserInfo gets user info from DB
func GetUserInfo(DB *sql.DB, login string) (data.User, error) {
	var user data.User

	stmt := "SELECT * FROM profiles WHERE username = $1"
	err := DB.QueryRow(stmt, login).Scan(&user.UserID, &user.UserName, &user.UserFIO, &user.UserisAdmin, &user.UserRate)
	if err != nil {
		return user, err
	}

	return user, nil
}

//GetAllUsers gets all users from DB
func GetAllUsers(DB *sql.DB) ([]data.User, error) {
	var users []data.User

	stmt := "SELECT * FROM profiles"
	rows, err := DB.Query(stmt)
	if err != nil {
		return users, err
	}
	defer rows.Close()

	for rows.Next() {
		var user data.User
		err = rows.Scan(&user.UserID, &user.UserName, &user.UserFIO, &user.UserisAdmin, &user.UserRate)
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

func GetGroupUsers(DB *sql.DB, groupName string) ([]data.User, error) {
	var (
		groupID int
		users []data.User
	)

	stmt := "SELECT id FROM groups WHERE group_name = $1"
	err := DB.QueryRow(stmt, groupName).Scan(&groupID)
	if err != nil {
		return users, err
	}

	rows, err := DB.Query("SELECT * FROM profiles WHERE id = (SELECT profile_id FROM groups_profiles WHERE group_id = $1", groupID)
	if err != nil {
		return users, err
	}
	defer rows.Close()

	for rows.Next() {
		var user data.User

		err = rows.Scan(&user.UserID, &user.UserName, &user.UserFIO, &user.UserisAdmin, &user.UserRate)
		if err != nil {
			return users, nil
		}

		users = append(users, user)
	}

	return users, nil
}

//GetTask gets info of task from DB
func GetTask(DB *sql.DB, ID int) (data.Task, error) {
	var (
		task      data.Task
		execID    int
		creatorID int
	)

	stmt := "SELECT * FROM tasks WHERE id = $1"
	err := DB.QueryRow(stmt, ID).Scan(&task.TaskID, &task.TaskName, &creatorID, &task.TaskExecutorType, &execID, &task.TaskStat, &task.TaskDateStart, &task.TaskDateEnd, &task.TaskRate)
	if err != nil {
		return task, err
	}

	task.TaskCreator, err = helpers.ConvertIDToExec(creatorID, "user")
	if err != nil {
		return task, err
	}

	task.TaskExecutor, err = helpers.ConvertIDToExec(execID, task.TaskExecutorType)
	if err != nil {
		return task, err
	}

	switch task.TaskExecutorType {
	case "user":
		{
			err = DB.QueryRow("SELECT username FROM profiles WHERE user_fio = $1", task.TaskExecutor).Scan(&task.TaskExecutorName)
			if err != nil {
				return task, err
			}
		}
	case "group":
		{
			task.TaskExecutorName = task.TaskExecutor
		}
	}

	err = DB.QueryRow("SELECT username FROM profiles WHERE user_fio = $1", task.TaskCreator).Scan(&task.TaskCreatorName)
	if err != nil {
		return task, err
	}

	return task, nil
}

//GetAllTasks gets all tasks from DB
func GetAllTasks(DB *sql.DB) ([]data.Task, error) {
	var tasks []data.Task

	stmt := "SELECT * FROM tasks"
	rows, err := DB.Query(stmt)
	if err != nil {
		return tasks, err
	}
	defer rows.Close()

	for rows.Next() {
		var (
			task      data.Task
			creatorID int
			execID    int
		)
		err = rows.Scan(&task.TaskID, &task.TaskName, &creatorID, &task.TaskExecutorType, &execID, &task.TaskStat, &task.TaskDateStart, &task.TaskDateEnd, &task.TaskRate)
		if err != nil {
			return tasks, err
		}

		task.TaskExecutor, err = helpers.ConvertIDToExec(execID, task.TaskExecutorType)
		if err != nil {
			return tasks, err
		}

		task.TaskCreator, err = helpers.ConvertIDToExec(creatorID, "user")
		if err != nil {
			return tasks, err
		}

		switch task.TaskExecutorType {
		case "user":
			{
				err = DB.QueryRow("SELECT username FROM profiles WHERE user_fio = $1", task.TaskExecutor).Scan(&task.TaskExecutorName)
				if err != nil {
					return tasks, err
				}
			}
		case "group":
			{
				task.TaskExecutorName = task.TaskExecutor
			}
		}

		err = DB.QueryRow("SELECT username FROM profiles WHERE user_fio = $1", task.TaskCreator).Scan(&task.TaskCreatorName)
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

//GetAllGroups gets all users from DB
func GetAllGroups(DB *sql.DB) ([]data.Group, error) {
	var groups []data.Group

	stmt := "SELECT * FROM groups"
	rows, err := DB.Query(stmt)
	if err != nil {
		return groups, err
	}
	defer rows.Close()

	for rows.Next() {
		var group data.Group
		err = rows.Scan(&group.GroupID, &group.GroupName)
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
