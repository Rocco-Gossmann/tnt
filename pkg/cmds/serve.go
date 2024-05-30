package cmds

import (
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"runtime"
	"time"

	"github.com/rocco-gossmann/tnt/pkg/serve"
	"github.com/rocco-gossmann/tnt/pkg/utils"

	"github.com/spf13/cobra"
)

type HandlerFunc func(writer http.ResponseWriter, request *http.Request)

var openCMD = ""

func globalHeaders(next HandlerFunc) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("pragma", "no-cache")
		w.Header().Set("cache-control", "post-check=0, pre-check=0, no-store, no-cache, must-revalidate, max-age=0")
		w.Header().Set("expires", "Thu, 01 Jan 1970 00:00:00 GMT")

		next(w, r)
	}
}

func logRequestPrefix(prefix string, next HandlerFunc) HandlerFunc {
	log.SetPrefix(prefix)
	return logRequest(next)
}

func logRequest(next HandlerFunc) HandlerFunc {
	prefix := log.Prefix()
	return func(w http.ResponseWriter, r *http.Request) {
		log.SetPrefix(fmt.Sprintf("%s -> [%s] %s => ", prefix, r.Method, r.URL.Path))
		log.Println("called")
		next(w, r)
	}
}

func runOpen(cmd string, url string) {
	log.Println("let's goooooooo !!!!!!")
	dur, err := time.ParseDuration("1s")
	if err != nil {
		fmt.Println("[Warning!] Can't open browser", err)
		return
	}
	log.Println("waiting for:", dur)
	t := time.NewTimer(dur)
	<-t.C

	log.Println("run xdg-open", dur)

	err = exec.Command(cmd, url).Run()
	if err != nil {
		fmt.Println("[Warning!] Can't open browser", err)
	}
	return
}

var ServeCMD cobra.Command = cobra.Command{
	Use: "serve [-p|--port]",
	Run: func(cmd *cobra.Command, args []string) {

		port, err := cmd.Flags().GetUint32("port")
		if err != nil {
			utils.Failf("serve error: %s", err)
		}

		url, err := cmd.Flags().GetString("url")
		if err != nil {
			utils.Failf("serve error: %s", err)
		}

		mux := http.NewServeMux()

		mux.HandleFunc("GET /", logRequestPrefix("(GET /)", globalHeaders(serve.GetIndex)))
		mux.HandleFunc("GET /htmx.js", logRequestPrefix("(GET /htmx.js)", serve.FileServer("views/htmx.js")))
		mux.HandleFunc("GET /main.css", logRequestPrefix("(GET /main.css)", globalHeaders(serve.FileServer("views/main.css"))))

		mux.HandleFunc("POST /task", logRequestPrefix("(POST /task)", globalHeaders(serve.PostTask)))
		mux.HandleFunc("GET /tasks", logRequestPrefix("(GET /tasks)", globalHeaders(serve.GetTasks)))
		mux.HandleFunc("DELETE /task/{id}", logRequestPrefix("(DELETE /task/{id})", globalHeaders(serve.DeleteTask)))

		mux.HandleFunc("POST /timer/{taskid}", logRequestPrefix("(POST /timer/{taskid})", globalHeaders(serve.PostTime)))
		mux.HandleFunc("DELETE /timer", logRequestPrefix("(DELETE /timer)", globalHeaders(serve.EndTime)))
		mux.HandleFunc("DELETE /timer/{taskid}", logRequestPrefix("(DELETE /timer/{taskid})", globalHeaders(serve.EndTime)))

		mux.HandleFunc("GET /times", logRequestPrefix("(GET /times)", globalHeaders(serve.GetTimes)))
		mux.HandleFunc("GET /times/{taskid}", logRequestPrefix("(GET /times/{taskid})", globalHeaders(serve.GetTimes)))
		mux.HandleFunc("GET /times/sum", logRequestPrefix("(GET /times/sum)", globalHeaders(serve.GetTimeSums)))
		mux.HandleFunc("GET /times/sum/{taskid}", logRequestPrefix("(GET /times/sum/{taskid})", globalHeaders(serve.GetTimeSums)))

		mux.HandleFunc("DELETE /time/{timeid}", logRequestPrefix("(DELETE /time/{timeid})", globalHeaders(serve.DeleteTime)))

		mux.HandleFunc("GET /time/edit/{timeid}", logRequestPrefix("(GET /time/edit/{timeid})", globalHeaders(serve.GetTimeEdit)))

		addr := fmt.Sprintf("%s:%d", url, port)

		server := http.Server{
			Addr:    addr,
			Handler: mux,
		}

		bDoOpen, err := cmd.Flags().GetBool("open")
		if err != nil {
			fmt.Println("[Warning] failed to check --open flag", err)
		} else if bDoOpen {
			switch {
			case runtime.GOOS == "macos":
				go runOpen("open", addr)
			case runtime.GOOS == "linux":
				go runOpen("xdg-open", addr)
			default:
				fmt.Println("can't open on '", runtime.GOOS, "' yet") // TODO: Implement
			}
		}

		fmt.Printf("now listening at: %s", server.Addr)
		err = server.ListenAndServe()
		if err != nil {
			utils.Failf("serve error: %s", err)
		}
	},
}

func init() {
	ServeCMD.Flags().Uint32P("port", "p", 7353, "the port on which to serve the Web-Interface")
	ServeCMD.Flags().StringP("url", "u", "0.0.0.0", "the host adress on which to host the UI-Server")
	ServeCMD.Flags().BoolP("open", "o", false, "try to open the interface in your browser")

}
