package serve

import (
	"log"
	"net/http"
	"strconv"

	"github.com/rocco-gossmann/tnt/pkg/database"
)

func GetTimes(w http.ResponseWriter, r *http.Request) {

	var iTaskID uint = 0

	log.SetPrefix("GET /times => ")
	log.Println("called GET /times ")

	noCacheHeaders(&w)

	sTaskID := r.PathValue("taskid")

	if sTaskID != "" {
		tmp, err := strconv.ParseInt(sTaskID, 10, 64)
		if serveErr(&w, err) {
			return
		}

		iTaskID = uint(tmp)
	}

	times, err := database.GetTimesRaw(iTaskID)
	if serveErr(&w, err) {
		return
	}

	if times == nil {
		serveStatusMsg(&w, http.StatusNoContent, "no tasks")
	} else {
		for _, time := range times {
			log.Println("Render:", time);
			tmpl.ExecuteTemplate(w, "times_list", time)
		}
	}

}
