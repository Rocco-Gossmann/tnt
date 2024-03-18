package cmds

import (
	"log"
	"net/http"

	"github.com/rocco-gossmann/tnt/pkg/serve"

	"github.com/spf13/cobra"
)

func FileServer(fl string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) { http.ServeFile(w, r, fl) }
}

var ServeCMD cobra.Command = cobra.Command{
	Use: "serve [-p|--port]",
	Run: func(cmd *cobra.Command, args []string) {

		log.SetPrefix("tnt-server")

		mux := http.NewServeMux()

		mux.HandleFunc("GET /", serve.GetIndex)
		mux.HandleFunc("GET /htmx.js", FileServer("views/htmx.js"))

		mux.HandleFunc("POST /task", serve.PostTask)
		mux.HandleFunc("GET /tasks", serve.GetTasks)
		mux.HandleFunc("DELETE /task/{id}", serve.DeleteTask)

		mux.HandleFunc("GET /times/{taskid}", serve.GetTimes)
		mux.HandleFunc("GET /times", serve.GetTimes)

		server := http.Server{
			Addr: "0.0.0.0:7353",
			Handler: mux,
		}


		server.ListenAndServe()

	},
}
