package database

import (
	"database/sql"
	"errors"
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
	Id       uint
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
				unixepoch(ti.end) - unixepoch(ti.start) duration,
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
		id, tid  uint
		s, e, n  sql.NullString
		d        string
		duration sql.NullFloat64
	)

	for res.Next() {

		err := res.Scan(&id, &s, &e, &duration, &tid, &n)

		log.Println(id, s, e, d, tid, n)

		if !s.Valid {
			s.String = ""
		}

		if !e.Valid {
			e.String = "** running **"
		}

		if !duration.Valid {
			duration.Float64 = 0.00
		}

		if !n.Valid {
			n.String = " unknown "
		}

		if err != nil {
			return nil, err
		}

		d = utils.SecToTimePrint(duration.Float64)

		ret = append(ret, Time{
			Id:       id,
			Start:    s.String,
			End:      e.String,
			Duration: d,
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
			ORDER BY start ASC;
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

func DeleteTime(iTimeId uint) error {
	_, err := ExecStatement("DELETE FROM times WHERE id=?", iTimeId)
	return err
}

func GetTimeByID(iTimeId uint) (ds TimeDS, err error) {

	str := fmt.Sprintf(`
			SELECT
				ti.id,
				ta.name,
				ti.start,
				ti.end,
				unixepoch(ti.end) - unixepoch(ti.start) duration
			FROM times ti
			LEFT JOIN tasks ta ON ti.taskId = ta.id
			WHERE ti.id=%d
			ORDER BY start ASC;
		`, iTimeId)

	log.Print("query: ", str)

	res, err := QueryStatement(str)
	if err != nil {
		return
	}
	defer res.Close()

	if res.Next() {
		var name, total, start, end sql.NullString
		var id sql.NullInt64

		err = res.Scan(&id, &name, &start, &end, &total)
		utils.Err(err)

		ds.Id = iTimeId
		ds.Task = name.String
		ds.Start = start.String
		ds.End = end.String
		ds.Duration = total.String

	} else {
		err = errors.New("time not found")

	}

	return
}

func UpdateTimeDataset(iTimeId uint, start time.Time, end time.Time) (err error) {

	sSQLStart := start.Format("2006-01-02 15:04:05")
	sSQLEnd := end.Format("2005-01-02 15:04:05")

	_, err = ExecStatement(
		`UPDATE times SET start = ?, end = ? WHERE id = ?`,
		sSQLStart, sSQLEnd,
		iTimeId,
	)

	return
}
