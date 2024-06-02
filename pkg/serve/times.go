package serve

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/rocco-gossmann/tnt/pkg/database"
	"github.com/rocco-gossmann/tnt/pkg/utils"
)

// BM: ServeFunctions
// ==============================================================================

// BM: POST /timer/{taskId}
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

// BM: DELETE /timer
func EndTime(w http.ResponseWriter, r *http.Request) {
	runInit()

	database.FinishCurrentlyRunningTimes()

	GetTasks(w, r)
	GetTimes(w, r)
}

// BM: GET /times
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

// BM: DELETE /time/{timeId}
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

// BM: GET /times/sum   &&   GET /times/sum/{taskId}
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

// BM: GET /time/edit/{timeid}
func GetTimeEdit(w http.ResponseWriter, r *http.Request) {

	_, oTimeEdit, err := getEditTimeByPathValue(r, "timeid")
	if serveErr(&w, err) {
		return
	}

	tmpl.ExecuteTemplate(w, "time_edit_row", oTimeEdit)
}

// BM: POST /time/edit/{timeid}
func PostTimeEdit(w http.ResponseWriter, r *http.Request) {

	iTimeId, oTimeEdit, err := getEditTimeByPathValue(r, "timeid")
	if serveErr(&w, err) {
		return
	}

	err = r.ParseForm()
	if serveErr(&w, err) {
		return
	}

	body := r.Form
	log.Println("hit", body)

	oTimeEdit.StartDate = readTimeEditFieldFromBody(body, "startdate", oTimeEdit.StartDate)
	oTimeEdit.StartTime = readTimeEditFieldFromBody(body, "starttime", oTimeEdit.StartTime)
	oTimeEdit.EndDate = readTimeEditFieldFromBody(body, "enddate", oTimeEdit.EndDate)
	oTimeEdit.EndTime = readTimeEditFieldFromBody(body, "endtime", oTimeEdit.EndTime)

	oStartTime, err := time.Parse("2006-01-02 15:04", fmt.Sprintf("%s %s", oTimeEdit.StartDate, oTimeEdit.StartTime))
	if serveErr(&w, err) {
		return
	}

	oEndTime, err := time.Parse("2006-01-02 15:04", fmt.Sprintf("%s %s", oTimeEdit.EndDate, oTimeEdit.EndTime))
	if serveErr(&w, err) {
		return
	}

	err = database.UpdateTimeDataset(uint(iTimeId), oStartTime, oEndTime)
	if serveErr(&w, err) {
		return
	}

	respondWithTime(&w, uint(iTimeId))
}

// BM: GET /time/{timeid}
func GetTime(w http.ResponseWriter, r *http.Request) {
	iTimeId, _, err := getEditTimeByPathValue(r, "timeid")
	if serveErr(&w, err) {
		return
	}

	respondWithTime(&w, uint(iTimeId))
}

// BM: Helper Functions
// ==============================================================================
type cTimeEditTime struct {
	Id                   uint
	Task                 string
	Duration             string
	StartDate, StartTime string
	EndDate, EndTime     string
}

func readTimeEditFieldFromBody(u url.Values, field string, fallback string) string {
	if u.Has(field) {
		return u.Get(field)
	} else {
		return fallback
	}
}

func respondWithTime(w *http.ResponseWriter, iTimeId uint) {

	oTime, err := database.GetTimeByIDRaw(iTimeId)
	if serveErr(w, err) {
		return
	}

	tmpl.ExecuteTemplate(*w, "time_list_entry", oTime)
}

func getEditTimeByPathValue(r *http.Request, pathValue string) (iTimeID int64, oTimeEdit cTimeEditTime, err error) {

	var oTime database.TimeDS
	iTimeID, err = strconv.ParseInt(r.PathValue(pathValue), 10, 64)
	if err == nil {
		oTime, err = database.GetTimeByID(uint(iTimeID))
		if err != nil {
			return
		}
	} else {
		return
	}

	// Init TimeEdit
	//==============================================================================
	oTimeEdit.Id = oTime.Id
	oTimeEdit.Task = oTime.Task

	// Duration
	//==============================================================================
	iDur, err := strconv.ParseFloat(oTime.Duration, 64)
	if err != nil {
		return
	}
	oTimeEdit.Duration = utils.SecToTimePrint(iDur)

	// StartTime
	//==============================================================================
	t, err := time.Parse(utils.SQL_OUTPUT_DATETIMEFORMAT, oTime.Start)
	if err != nil {
		return
	}
	oTimeEdit.StartDate = t.Format("2006-01-02")
	oTimeEdit.StartTime = t.Format("15:04")

	// EndTime
	//==============================================================================
	t, err = time.Parse(utils.SQL_OUTPUT_DATETIMEFORMAT, oTime.End)
	if err != nil {
		return
	}
	oTimeEdit.EndDate = t.Format("2006-01-02")
	oTimeEdit.EndTime = t.Format("15:04")

	return
}
