package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/rocco-gossmann/tnt/pkg/utils"
)

func initTimesTable() {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS times (
			id INTEGER PRIMARY KEY, 
			taskId INTEGER, 
			start DATETIME,
			end DATETIME,
			FOREIGN KEY(taskId) REFERENCES tasks(id)
		)
	`)

	utils.Err(err)
}

type TimeDS struct {
	task     string
	start    string
	end      string
	duration string
}

func (ts TimeDS) String() string {
	return fmt.Sprintf(" %s | %s | %s | %s ", ts.task, ts.start, ts.end, ts.duration)
}

func getSQLTimeNow() string {
	tm := time.Now()
	return tm.Format(utils.SQL_DATETIME_FORMAT)
}

func TimedTaskIsRunning(taskId uint) bool {
	c := 0
	result, err := RowQueryStatement("SELECT COUNT(*) FROM times WHERE taskId=? AND end IS NULL", taskId)
	utils.Err(err)

	err = result.Scan(&c)
	utils.Err(err)

	return c > 0
}

// @panic error - on DB issue
func FinishCurrentlyRunningTimes() {
	endTS := getSQLTimeNow()
	result, err := ExecStatement("UPDATE times SET end=? WHERE END IS NULL", endTS)
	utils.Err(err)

	rowCnt, err := result.RowsAffected()
	utils.Err(err)

	suffix := utils.Suffix(int(rowCnt), "y", "ies")

	log.Printf("Ended %d running `time` entr%s", rowCnt, suffix)
}

// @panic error - on db fail
func StartNewTime(taskId uint) int {
	startTS := getSQLTimeNow()

	result, err := ExecStatement("INSERT INTO times(taskId, start) values(?, ?)", taskId, startTS)
	utils.Err(err)

	insertId, err := result.LastInsertId()
	utils.Err(err)

	return int(insertId)

}

func GetTimes(taskId uint) ([]TimeDS, error) {
	var ret []TimeDS
	var taskWhere = ""

	if taskId > 0 {
		taskWhere = fmt.Sprintf(" WHERE taskId=%d ", taskId)
	}

	res, err := QueryStatement(`
			SELECT 
				ta.name, 
				ti.start, 
				ti.end, 
				time(unixepoch(ti.end) - unixepoch(ti.start), "unixepoch") duration
			FROM times ti
			LEFT JOIN tasks ta ON ti.taskId = ta.id
			` + taskWhere + `
			ORDER BY start DESC;
		`)

	if err != nil {
		return ret, err
	}

	for res.Next() {
		var name, total, start, end sql.NullString
		err = res.Scan(&name, &start, &end, &total)
		utils.Err(err)

		if !end.Valid {
			end.String = "* running *"
		}

		if total.Valid {
			total.String += " Hours"
		}

		ret = append(ret, TimeDS{
			task:     name.String,
			start:    start.String,
			end:      end.String,
			duration: total.String,
		})
	}

	return ret, err
}
