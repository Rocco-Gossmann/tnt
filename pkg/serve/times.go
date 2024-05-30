package serve

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/rocco-gossmann/tnt/pkg/database"
	"github.com/rocco-gossmann/tnt/pkg/utils"
)

func PostTime(w http.ResponseWriter, r *http.Request) {
	runInit()

	log.Println("called: PostTime")

	iTaskID, err := strconv.ParseInt(r.PathValue("taskid"), 10, 64)
	if serveErr(&w, err) {
		return
	}
	// Must get Task object only after task is running (otherwise buttons are wrong)
	if database.TimedTaskIsRunning(uint(iTaskID)) {
		serveStatusMsg(&w, http.StatusAccepted, "&#x23F9;")
		return
	}
	database.FinishCurrentlyRunningTimes()

	_, err = database.StartNewTimeRaw(uint(iTaskID))
	if serveErr(&w, err) {
		return
	}

	GetTasks(w, r)
	GetTimes(w, r)
}

func EndTime(w http.ResponseWriter, r *http.Request) {
	runInit()

	database.FinishCurrentlyRunningTimes()

	GetTasks(w, r)
	GetTimes(w, r)
}

func GetTimes(w http.ResponseWriter, r *http.Request) {

	runInit()

	var task database.Task

	sTaskID := r.PathValue("taskid")

	if sTaskID != "" {
		tmp, err := strconv.ParseInt(sTaskID, 10, 64)
		if serveErr(&w, err) {
			return
		}

		task, err = database.GetTaskById(uint(tmp))
		if serveErr(&w, err) {
			return
		}
	}

	times, err := database.GetTimesRaw(task.Id)
	if serveErr(&w, err) {
		return
	}

	tmpl.ExecuteTemplate(w, "times_list_section", times)
	tmpl.ExecuteTemplate(w, "times_list_label", task.Name)
}

func DeleteTime(w http.ResponseWriter, r *http.Request) {
	runInit()

	iTimeID, err := strconv.ParseInt(r.PathValue("timeid"), 10, 64)
	if serveErr(&w, err) {
		return
	}

	err = database.DeleteTime(uint(iTimeID))
	if serveErr(&w, err) {
		return
	}

	GetTimes(w, r)
}

func GetTimeSums(w http.ResponseWriter, r *http.Request) {

	var (
		iTaskID int64 = 0
		err     error
	)

	sTaskId := r.PathValue("taskid")
	if len(sTaskId) > 0 {
		iTaskID, err = strconv.ParseInt(sTaskId, 10, 64)
		utils.Err(err)
	}

	sums := database.GetTimeSums(uint(iTaskID))

	tmpl.ExecuteTemplate(w, "times_sum_label", nil)
	tmpl.ExecuteTemplate(w, "times_sum_section", sums)
}

type cTimeEditTime struct {
	Id                   uint
	Task                 string
	Duration             string
	StartDate, StartTime string
	EndDate, EndTime     string
}

func GetTimeEdit(w http.ResponseWriter, r *http.Request) {
	iTimeID, err := strconv.ParseInt(r.PathValue("timeid"), 10, 64)
	if serveErr(&w, err) {
		return
	}

	oTime, err := database.GetTimeByID(uint(iTimeID))
	if serveErr(&w, err) {
		return
	}

	// Init TimeEdit
	//==============================================================================
	var oTimeEdit cTimeEditTime
	oTimeEdit.Id = oTime.Id
	oTimeEdit.Task = oTime.Task

	// Duration
	//==============================================================================
	iDur, err := strconv.ParseFloat(oTime.Duration, 10)
	if serveErr(&w, err) {
		return
	}
	oTimeEdit.Duration = utils.SecToTimePrint(iDur)

	// StartTime
	//==============================================================================
	t, err := time.Parse(utils.SQL_OUTPUT_DATETIMEFORMAT, oTime.Start)
	if serveErr(&w, err) {
		return
	}
	oTimeEdit.StartDate = t.Format("2006-01-02")
	oTimeEdit.StartTime = t.Format("15:04:05")

	// EndTime
	//==============================================================================
	t, err = time.Parse(utils.SQL_OUTPUT_DATETIMEFORMAT, oTime.End)
	if serveErr(&w, err) {
		return
	}
	oTimeEdit.EndDate = t.Format("2006-01-02")
	oTimeEdit.EndTime = t.Format("15:04:05")

	tmpl.ExecuteTemplate(w, "time_edit_row", oTimeEdit)

}
