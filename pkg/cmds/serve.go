package cmds

import (
	"log"
	"net/http"

	"github.com/rocco-gossmann/tnt/pkg/serve"

	"github.com/spf13/cobra"
)

type HandlerFunc func(writer http.ResponseWriter, request *http.Request)

func globalHeaders(next HandlerFunc) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("pragma", "no-cache")
		w.Header().Set("cache-control", "post-check=0, pre-check=0, no-store, no-cache, must-revalidate, max-age=0")
		w.Header().Set("expires", "Thu, 01 Jan 1970 00:00:00 GMT")

		next(w, r)
	}
}

func FileServer(fl string) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) { http.ServeFile(w, r, fl) }
}

var ServeCMD cobra.Command = cobra.Command{
	Use: "serve [-p|--port]",
	Run: func(cmd *cobra.Command, args []string) {

		//database.DeInitDB()

		log.SetPrefix("tnt-server")

		mux := http.NewServeMux()

		mux.HandleFunc("GET /", globalHeaders(serve.GetIndex))
		mux.HandleFunc("GET /htmx.js", FileServer("views/htmx.js"))
		mux.HandleFunc("GET /main.css", globalHeaders(FileServer("views/main.css")))

		mux.HandleFunc("POST /task", globalHeaders(serve.PostTask))
		mux.HandleFunc("GET /tasks", globalHeaders(serve.GetTasks))
		mux.HandleFunc("DELETE /task/{id}", globalHeaders(serve.DeleteTask))

		mux.HandleFunc("POST /timer/{taskid}", globalHeaders(serve.PostTime))
		mux.HandleFunc("DELETE /timer", globalHeaders(serve.EndTime))
		mux.HandleFunc("DELETE /timer/{taskid}", globalHeaders(serve.EndTime))
		mux.HandleFunc("GET /times/{taskid}", globalHeaders(serve.GetTimes))
		mux.HandleFunc("DELETE /time/{timeid}", globalHeaders(serve.DeleteTime))
		mux.HandleFunc("GET /times", globalHeaders(serve.GetTimes))

		server := http.Server{
			Addr:    "0.0.0.0:7353",
			Handler: mux,
		}

		server.ListenAndServe()

	},
}
