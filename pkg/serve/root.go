package serve

import (
	"log"
	"net/http"
	"text/template"

	"github.com/rocco-gossmann/tnt/pkg/database"
)

var tmpl *template.Template

type IndexContext struct {
	Cnt uint
}

var context IndexContext = IndexContext{
	Cnt: 0,
}

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
		t, err := template.ParseFiles("views/index.html")

		if err != nil {
			log.Fatal("failed to parse files", err)
		}

		tmpl = t
	}

	database.InitDB("")
}

func runDeInit() {
	if tmpl != nil {
		tmpl = nil
	}

	database.DeInitDB()
}
