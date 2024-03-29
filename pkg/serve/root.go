package serve

import (
	"log"
	"net/http"
	"text/template"
)


var tmpl *template.Template

func serveErr(res *http.ResponseWriter, err error) bool {

	if err != nil {
		log.Println("Encountered error: ", err)
		(*res).WriteHeader(http.StatusInternalServerError)
		(*res).Write([]byte(err.Error()))
		return true
	}

	return false
}

func serveStatusMsg(w *http.ResponseWriter, status int, msg string) {
	(*w).WriteHeader(status)
	(*w).Write([]byte(msg))
	return
}

func makeResponseJSON(res *http.ResponseWriter) {
	(*res).Header().Set("content-type", "application/json")
}

func runInit() {
	if tmpl == nil {
		t, err := template.ParseFS(views, "views/index.html")

		if err != nil {
			log.Fatal("failed to parse files", err)
		}

		tmpl = t
	}
}

func FileServer(fl string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) { http.ServeFileFS(w, r, views, fl) }
}

