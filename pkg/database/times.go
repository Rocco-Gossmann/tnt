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

type Time struct {
	Id       uint
	Start    string
	End      string
	Duration string
	TaskId   uint
	TaskName string
}

type TimeDS struct {
	Task     string
	Start    string
	End      string
	Duration string
}

func (ts TimeDS) String() string {
	return fmt.Sprintf(" %s | %s | %s | %s ", ts.Task, ts.Start, ts.End, ts.Duration)
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

func StartNewTimeRaw(taskId uint) (insertId int64, err error) {
	startTS := getSQLTimeNow()

	result, err := ExecStatement("INSERT INTO times(taskId, start) values(?, ?)", taskId, startTS)
	if err != nil {
		return
	}

	insertId, err = result.LastInsertId()

	return
}

// @panic error - on db fail
func StartNewTime(taskId uint) int {

	insertId, err := StartNewTimeRaw(taskId)
	utils.Err(err)

	return int(insertId)
}

func GetTimesRaw(taskId uint) ([]Time, error) {

	var ret []Time

	var taskWhere = ""

	if taskId > 0 {
		taskWhere = fmt.Sprintf(" WHERE taskId=%d ", taskId)
	}

	res, err := QueryStatement(`
			SELECT 
				ti.id, 
				ti.start, 
				ti.end, 
				time(unixepoch(ti.end) - unixepoch(ti.start), "unixepoch") duration,
				ti.taskId,
				ta.name
			FROM times ti
			LEFT JOIN tasks ta ON ti.taskId = ta.id
			` + taskWhere + `
			ORDER BY start DESC;
		`)

	if err != nil {
		return nil, err
	}

	var (
		id, tid    uint
		s, e, d, n sql.NullString
	)

	for res.Next() {

		err := res.Scan(&id, &s, &e, &d, &tid, &n)

		log.Println(id, s, e, d, tid, n)

		if !s.Valid {
			s.String = ""
		}

		if !e.Valid {
			e.String = "** running **"
		}

		if !d.Valid {
			d.String = ""
		}

		if !n.Valid {
			n.String = " unknown "
		}

		if err != nil {
			return nil, err
		}

		ret = append(ret, Time{
			Id:       id,
			Start:    s.String,
			End:      e.String,
			Duration: d.String,
			TaskId:   tid,
			TaskName: n.String,
		})
	}

	return ret, nil
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
			Task:     name.String,
			Start:    start.String,
			End:      end.String,
			Duration: total.String,
		})
	}

	return ret, err
}
