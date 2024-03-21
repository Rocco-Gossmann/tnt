package serve

import (
	"log"
	"net/http"
	"strconv"

	"github.com/rocco-gossmann/tnt/pkg/database"
)

func GetTimes(w http.ResponseWriter, r *http.Request) {

	runInit()
	defer runDeInit()

	var task database.Task

	log.SetPrefix("GET /times => ")
	log.Println("called GET /times ")

	noCacheHeaders(&w)

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

	if times == nil {
		serveStatusMsg(&w, http.StatusNoContent, "no tasks")
	} else {
		tmpl.ExecuteTemplate(w, "times_section", context)
	}

}
