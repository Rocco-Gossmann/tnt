package serve

import (
	"log"
	"net/http"
	"strconv"

	"github.com/rocco-gossmann/tnt/pkg/database"
)

func GetTasks(w http.ResponseWriter, r *http.Request) {

	runInit()
	defer runDeInit()

	log.SetPrefix("GET /tasks => ")
	log.Println("called GET /tasks ")
	noCacheHeaders(&w)

	tasks, err := database.GetTaskList()
	if serveErr(&w, err) {
		return
	}

	for _, t := range tasks {
		tmpl.ExecuteTemplate(w, "task_list", t)
	}

}

func PostTask(w http.ResponseWriter, r *http.Request) {

	runInit()
	defer runDeInit()

	log.SetPrefix("POST /task => ")
	log.Println("called POST /task ", r)

	err := r.ParseForm()
	serveErr(&w, err)

	log.Println("Form: ", r.PostForm)

	if !r.PostForm.Has("taskname") {
		serveStatusMsg(&w, http.StatusBadRequest, "missing taskname")
		return
	}

	taskName := r.PostForm.Get("taskname")
	err = database.AddTask(taskName)

	if database.IsUniqueContraintError(err) {
		serveStatusMsg(&w, http.StatusConflict, "task already exists")
		return

	} else if serveErr(&w, err) {
		return

	}

	t, err := database.GetTaskByName(taskName)
	if !serveErr(&w, err) {
		tmpl.ExecuteTemplate(w, "task_list", t)
	}

}

func DeleteTask(w http.ResponseWriter, r *http.Request) {
	runInit()
	defer runDeInit()

	log.SetPrefix("DELETE /tasks => ")
	noCacheHeaders(&w)

	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	serveErr(&w, err)

	log.Printf("called DELETE /task/%d\n", id)

	if id > 0 {
		rows, err := database.DropTask(uint(id))
		serveErr(&w, err)

		if rows > 0 {
			log.Println("deleted rows:", rows)
			w.WriteHeader(http.StatusOK)

		} else {
			log.Println("did not delete any rows")
			w.WriteHeader(http.StatusNoContent)

		}

	} else {
		log.Printf("received id <= 0 (got: %d) => skip deletion\n", id)
		w.WriteHeader(http.StatusNoContent)
	}

}
