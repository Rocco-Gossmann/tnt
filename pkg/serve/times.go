package serve

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/rocco-gossmann/tnt/pkg/database"
)

func GetTimes(w http.ResponseWriter, r *http.Request) {
	log.SetPrefix("GET /times => ")
	log.Println("called GET /times ")

	noCacheHeaders(&w)

	times, err := database.GetTimes(0)
	if serveErr(&w, err) {
		return
	}

	mar, err := json.MarshalIndent(times[0], "", "\t")
	if serveErr(&w, err) {
		return
	}

	makeResponseJSON(&w)
	w.Write(mar)
}
