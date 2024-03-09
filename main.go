package main

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/rocco-gossmann/tnt/pkg/cmds"
	"github.com/rocco-gossmann/tnt/pkg/cmds/tasks"
	"github.com/rocco-gossmann/tnt/pkg/cmds/times"
	"github.com/rocco-gossmann/tnt/pkg/database"
	"github.com/rocco-gossmann/tnt/pkg/utils"
	"github.com/spf13/cobra"

	ex "github.com/rocco-gossmann/go_throwable"
)

// main
func main() {

	hadAPanic := false

	//	log.SetOutput(io.Discard)
	// TO make sure the DB Connection can be closed safely, we need ot use panics
	// (os.Exit is bad)
	// but we don't want to overwhelm the user with some uggly messages, when a panic
	// We'll print a nice looking error instead
	// was caused on purpose. So if a Pnaic was caused by a utils.ControlledPanic,
	//
	// the go_throwable (ex) package helps us to achieve this easy
	ex.Try(func() any {

		myCMD := cobra.Command{
			Use: "tnt {tasks|start|switch|times}",
			PersistentPreRun: func(cmd *cobra.Command, args []string) {
				if cmd.Flag("debug").Value.String() == "false" {
					log.SetOutput(io.Discard)
					log.Println("--debug set => enable logging")
				}

				database.InitDB("")
			},
		}

		myCMD.AddCommand(&tasks.TaskCMD, &cmds.SwitchCMD, &cmds.StopCmd, &times.TimesCMD)
		myCMD.PersistentFlags().Bool("debug", false, "Enable Debug-Log output")

		myCMD.Execute()

		return nil
	}, ex.TryOpts{
		SkipWarnings: true,
		Catch: func(panicCause any) any {

			if ce, ok := interface{}(panicCause).(utils.ControlledPanic); ok {
				fmt.Println(ce.Msg)
				hadAPanic = true
			} else {
				fmt.Println("PANIC !!!! *** Runs in Circles ***")
				panic(panicCause)

			}
			return nil
		},
		Finally: func() {
			database.DeInitDB()
			if hadAPanic {
				os.Exit(1)
			}
		},
	})

}
