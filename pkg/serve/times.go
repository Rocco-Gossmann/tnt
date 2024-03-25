package serve

import (
	"net/http"
	"strconv"

	"github.com/rocco-gossmann/tnt/pkg/database"
)

func PostTime(w http.ResponseWriter, r *http.Request) {
	runInit()

	iTaskID, err := strconv.ParseInt(r.PathValue("taskid"), 10, 64)
	if serveErr(&w, err) {
		return
	}

	if database.TimedTaskIsRunning(uint(iTaskID)) {
		serveStatusMsg(&w, http.StatusAccepted, "&#x23F9;")
		return
	}
	database.FinishCurrentlyRunningTimes()

	_, err = database.StartNewTimeRaw(uint(iTaskID))
	if serveErr(&w, err) {
		return
	}

	serveStatusMsg(&w, http.StatusCreated, "OK")
}

func EndTime(w http.ResponseWriter, r *http.Request) {
	runInit()

	database.FinishCurrentlyRunningTimes()

	serveStatusMsg(&w, http.StatusNoContent, "OK")
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

	context := struct {
		Label string
		Times []database.Time
	}{
		Label: task.Name,
		Times: times,
	}

	tmpl.ExecuteTemplate(w, "times_section", context)

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

	serveStatusMsg(&w, http.StatusOK, "OK")
}
