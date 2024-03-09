package database

import (
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

func FinishCurrentlyRunningTimes() {
	endTS := getSQLTimeNow()
	result, err := ExecStatement("UPDATE times SET end=? WHERE END IS NULL", endTS)
	utils.Err(err)

	rowCnt, err := result.RowsAffected()
	utils.Err(err)

	suffix := utils.Suffix(int(rowCnt), "y", "ies")

	log.Printf("Ended %d running `time` entr%s", rowCnt, suffix)
}

func StartNewTime(taskId uint) int {
	startTS := getSQLTimeNow()

	result, err := ExecStatement("INSERT INTO times(taskId, start) values(?, ?)", taskId, startTS)
	utils.Err(err)

	insertId, err := result.LastInsertId()
	utils.Err(err)

	return int(insertId)

}
