package serve

import (
	"log"
	"net/http"
)

func GetIndex(w http.ResponseWriter, r *http.Request) {

	runInit()

	if r.URL.Path == "/" {
		log.Println("loading index")
		tmpl.ExecuteTemplate(w, "index", nil)
	} else {
		log.Println("not found")
		serveStatusMsg(&w, 404, "not found")
	}

}
