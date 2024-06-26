package serve

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/rocco-gossmann/tnt/pkg/database"
)

type TaskListTemplateContext struct {
	Tasks  []database.Task
	Search string
}

func GetTasks(w http.ResponseWriter, r *http.Request) {

	log.Println("call GetTasks")
	runInit()

	context := TaskListTemplateContext{Search: ""}

	var err error
	err = r.ParseForm()
	if serveErr(&w, err) {
		return
	}

	if r.Form.Has("task_search") {
		context.Search = r.Form.Get("task_search")
		context.Tasks, err = database.GetTaskList(context.Search)
	} else {
		context.Tasks, err = database.GetTaskList("")
	}

	if serveErr(&w, err) {
		return
	}

	tmpl.ExecuteTemplate(w, "task_list_section", context)

}

func PostTask(w http.ResponseWriter, r *http.Request) {

	runInit()

	err := r.ParseForm()
	serveErr(&w, err)

	if !r.PostForm.Has("taskname") {
		serveStatusMsg(&w, http.StatusBadRequest, "missing taskname")
		return
	}

	taskName := r.PostForm.Get("taskname")
	taskName = strings.TrimSpace(taskName)
	err = database.AddTask(taskName)

	if database.IsUniqueContraintError(err) {
		serveStatusMsg(&w, http.StatusConflict, "task already exists")
		return

	} else if serveErr(&w, err) {
		return

	}

	GetTasks(w, r)
}

func DeleteTask(w http.ResponseWriter, r *http.Request) {
	runInit()

	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	serveErr(&w, err)

	if id > 0 {
		rows, err := database.DropTask(uint(id))
		serveErr(&w, err)

		if rows > 0 {
			log.Println("deleted rows:", rows)

			GetTimes(w, r)

		} else {
			log.Println("did not delete any rows")
			w.WriteHeader(http.StatusNoContent)

		}

	} else {
		log.Printf("received id <= 0 (got: %d) => skip deletion\n", id)
		w.WriteHeader(http.StatusNoContent)
	}

}
