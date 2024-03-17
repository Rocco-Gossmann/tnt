package database

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/rocco-gossmann/tnt/pkg/utils"
)

func initTasksTable() {
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS tasks (id INTEGER PRIMARY KEY, textkey TEXT UNIQUE, name TEXT)")
	utils.Err(err)
}

type Task struct {
	Id uint
	Key string
	Name string
}

type TaskList []Task

// GetTaskList returns a slice of all Tasknames currently registered
func GetTaskList() ([]Task, error) {
	var ret []Task

	stmt, err := db.Prepare("SELECT id, textkey, name FROM tasks")
	if err != nil {
		return ret, err
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return ret, err
	}

	var sName, sKey string
	var iID uint


	for rows.Next() {
		rows.Scan(&iID, &sKey, &sName)
		ret = append(ret, Task{
			Id: iID,
			Name: sName,
			Key: sKey,
		})
	}
	defer rows.Close()

	return ret, nil
}

func internal_GenerateTaskKey(taskName string) (string, uint, error) {
	taskKey := GenerateTaskKey(taskName)
	var taskId sql.NullInt64

	taskRow, err := RowQueryStatement("SELECT id FROM tasks WHERE textkey = ?", taskKey)
	utils.Err(err)
	err = taskRow.Scan(&taskId)

	if !taskId.Valid {
		taskId.Int64 = 0
		return taskKey, uint(taskId.Int64), nil
	}

	return taskKey, uint(taskId.Int64), err
}

func GetTaskByName(taskName string) (*Task, error){
	_, taskId, err := internal_GenerateTaskKey(taskName)
	if err != nil {
		return nil, err
	}

	res, err := QueryStatement("SELECT id, textkey, name from tasks where id=?", taskId);
	if err != nil {
		return nil, err
	}

	defer res.Close()

	var (
		id uint
		name string
		textkey string
	)

	if res.Next() {
		err = res.Scan(&id, &textkey, &name)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, nil
	}

	return &Task{
		Id: id,
		Key: textkey,
		Name: name,
	}, nil
}

func GetTaskIDByName(taskName string) uint {

	_, taskId, err := internal_GenerateTaskKey(taskName)

	if err != nil {
		errStr := fmt.Sprintf("%s", err)
		if strings.HasPrefix(errStr, "sql: no rows") {
			utils.Failf("Task '%s' is not in the List. use tnt tasks add \"%s\" to add it.", taskName, taskName)
		} else {
			panic(err)
		}
	}

	utils.Err(err)

	log.Printf("TaskID is: %d", taskId)

	return taskId
}

func (tasks TaskList) ExtractTaskListNames() []string {
	lst := make([]string, len(tasks), len(tasks));

	for i, t := range tasks {
		lst[i] = t.Name;
	}

	return lst
}

func AddTask(taskName string) error {
	taskKey := GenerateTaskKey(taskName)
	_, err := ExecStatement("INSERT INTO tasks(textkey, name) VALUES (?, ?)", taskKey, taskName)

	if err != nil {
		return err
	}

	return nil
}

func RenameTask(taskId uint, newName string) (sql.Result, error) {
	taskKey, newTaskId, err := internal_GenerateTaskKey(newName)
	utils.Err(err)

	if newTaskId != 0 {
		utils.Failf("Task '%s' already exists. Please choose another new name.", newName)
	}

	return ExecStatement("UPDATE tasks SET textkey = ?, name = ? WHERE id = ?", taskKey, newName, taskId)
}

func DropTask(taskId uint) (int64, error)  {
	_, err := ExecStatement("DELETE FROM times WHERE taskId = ?", taskId)
	if err != nil {
		return 0, err 
	}
	r, err :=  ExecStatement("DELETE FROM tasks WHERE id = ?", taskId)
	if err != nil {
		return 0, err
	}

	rows, err := r.RowsAffected()	
	if err != nil {
		return 0, err
	}


	return rows, nil
}

// converts the given Task into a normaized version, to help avoid putting in mutlipe tasks
func GenerateTaskKey(taskName string) string {
	return strings.ToLower(strings.TrimSpace(taskName))
}
