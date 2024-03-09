package database

import (
	"fmt"
	"log"
	"strings"

	"github.com/rocco-gossmann/tnt/pkg/utils"
)

func initTasksTable() {
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS tasks (id INTEGER PRIMARY KEY, textkey TEXT UNIQUE, name TEXT)")
	utils.Err(err)
}

// GetTaskList returns a slice of all Tasknames currently registered
func GetTaskList() ([]string, error) {
	var ret []string

	stmt, err := db.Prepare("SELECT name FROM tasks")
	if err != nil {
		return ret, err
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return ret, err
	}

	var retStr string
	for rows.Next() {
		rows.Scan(&retStr)
		ret = append(ret, retStr)
	}
	defer rows.Close()

	return ret, nil
}

func GetTaskIDByName(taskName string) uint {
	taskKey := GenerateTaskKey(taskName)
	var taskId uint

	taskRow, err := RowQueryStatement("SELECT id FROM tasks WHERE textkey = ?", taskKey)
	utils.Err(err)

	err = taskRow.Scan(&taskId)
	if err != nil {
		errStr := fmt.Sprintf("%s", err)
		if strings.HasPrefix(errStr, "sql: no rows") {
			utils.Exitf("Task '%s' is not in the List. use tnt tasks add \"%s\" to add it.", taskName, taskName)
		} else {
			panic(err)
		}
	}

	utils.Err(err)

	log.Printf("TaskID is: %d", taskId)

	return taskId
}

// converts the given Task into a normaized version, to help avoid putting in mutlipe tasks
func GenerateTaskKey(taskName string) string {
	return strings.ToLower(strings.TrimSpace(taskName))
}
