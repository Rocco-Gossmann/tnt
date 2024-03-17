package serve

import (
	"log"
	"net/http"
	"text/template"
)

var tmpl *template.Template

type IndexContext struct {
	Cnt uint
}

var context IndexContext = IndexContext{
	Cnt: 0,
}

func serveErr(res *http.ResponseWriter, err error)  bool {

	if err != nil {
		log.Println("Encountered error: ", err);
		(*res).WriteHeader(http.StatusInternalServerError);
		(*res).Write([]byte(err.Error()))
		return true;
	}

	return false
}

func serveStatusMsg(w *http.ResponseWriter, status int, msg string) {
	(*w).WriteHeader(status);	
	(*w).Write([]byte(msg));
	return
}

func noCacheHeaders(res *http.ResponseWriter) {
	(*res).Header().Set("pragma", "no-cache")
	(*res).Header().Set("cache-control", "post-check=0, pre-check=0, no-store, no-cache, must-revalidate, max-age=0")
	(*res).Header().Set("expires", "Thu, 01 Jan 1970 00:00:00 GMT")
}

func makeResponseJSON(res *http.ResponseWriter){
	(*res).Header().Set("content-type", "application/json")
}

func runInit() {
	if tmpl == nil {
		t, err := template.ParseFiles("views/index.html");

		if err != nil {
			log.Fatal("failed to parse files", err);
		}

		tmpl = t
	}
}
