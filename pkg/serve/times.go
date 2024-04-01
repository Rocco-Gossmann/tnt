package serve

import (
	"log"
	"net/http"
	"strconv"

	"github.com/rocco-gossmann/tnt/pkg/database"
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
	serveStatusMsg(&w, http.StatusNotImplemented, "not implemented")
}
