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
	Id     uint
	Key    string
	Name   string
	Active bool
}

type TaskList []Task

// GetTaskList returns a slice of all Tasknames currently registered
func GetTaskList() (ret []Task, err error) {

	// Pretpare statements
	stmt, err := db.Prepare(`
		SELECT 
			ta.id, 
			ta.textkey, 
			ta.name,
			(SELECT COUNT(*) FROM times ti WHERE ti.taskId=ta.id AND ti.End IS NULL) active
		FROM 
			tasks ta
	`)
	if err != nil {
		return
	}
	defer stmt.Close()

	// Run Query
	rows, err := stmt.Query()
	if err != nil {
		return
	}
	defer rows.Close()

	// Process Rows
	var t Task
	var active int

	for rows.Next() {
		err = rows.Scan(&t.Id, &t.Key, &t.Name, &active)
		if err != nil {
			return
		}
		t.Active = active == 1
		ret = append(ret, t)
	}

	return
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

func GetTaskByName(taskName string) (t Task, err error) {

	_, taskId, err := internal_GenerateTaskKey(taskName)

	if err != nil {
		return
	}

	t, err = GetTaskById(taskId)
	if err != nil {
		// Rewrite error, since user did not asks for id, but name
		err = fmt.Errorf("no task for name '%s' found", taskName)
	}

	return
}

func GetTaskById(iTaskId uint) (t Task, err error) {

	res, err := QueryStatement("SELECT id, textkey, name from tasks where id=?", iTaskId)

	if err != nil {
		return
	}

	defer res.Close()

	if res.Next() {
		err = res.Scan(&t.Id, &t.Key, &t.Name)
	} else {
		err = fmt.Errorf("no task for id '%d' found", iTaskId)
	}

	return
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
	lst := make([]string, len(tasks), len(tasks))

	for i, t := range tasks {
		lst[i] = t.Name
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

func DropTask(taskId uint) (int64, error) {
	_, err := ExecStatement("DELETE FROM times WHERE taskId = ?", taskId)
	if err != nil {
		return 0, err
	}
	r, err := ExecStatement("DELETE FROM tasks WHERE id = ?", taskId)
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
