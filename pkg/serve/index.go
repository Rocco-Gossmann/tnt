package serve

import (
	"log"
	"net/http"
)



func GetIndex(w http.ResponseWriter, r *http.Request) {

	log.SetPrefix(" GET / => ")
	log.Println("called GET / ", context, tmpl )

	runInit();
	context.Cnt++;


	tmpl.ExecuteTemplate(w, "index", context)

}
