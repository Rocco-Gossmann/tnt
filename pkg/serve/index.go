package serve

import (
	"log"
	"net/http"

	"github.com/rocco-gossmann/tnt/pkg/database"
)

type IndexContext struct {
	Tasks []database.Task
	Times []database.Time
	Label string
}

func GetIndex(w http.ResponseWriter, r *http.Request) {

	runInit()

	var context = IndexContext{
		Times: make([]database.Time, 0),
		Tasks: make([]database.Task, 0),
		Label: "",
	}

	tasks, err := database.GetTaskList()
	if serveErr(&w, err) {
		return
	}
	context.Tasks = tasks

	times, err := database.GetTimesRaw(0)
	if serveErr(&w, err) {
		return
	}

	for i, oTime := range times {
		times[i], err = prepareTimeObjForOutput(oTime)
		if serveErr(&w, err) {
			return
		}
	}

	context.Times = times

	if r.URL.Path == "/" {
		log.Println("loading index")
		tmpl.ExecuteTemplate(w, "index", context)
	} else {
		log.Println("not found")
		serveStatusMsg(&w, 404, "not found")
	}

}
