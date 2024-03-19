package main

import (
	"fmt"
	"os"

	"github.com/rocco-gossmann/tnt/pkg/cmds"
	"github.com/rocco-gossmann/tnt/pkg/database"
	"github.com/rocco-gossmann/tnt/pkg/env"
	"github.com/rocco-gossmann/tnt/pkg/utils"

	ex "github.com/rocco-gossmann/go_throwable"
)

var (
	Version string
)

// main
func main() {

	hadAPanic := false
	panicExitCode := 1

	env.Version = Version + ""

	// To make sure the DB Connection can be closed safely, we need ot use panics
	// (os.Exit is bad)
	// but we don't want to overwhelm the user with some uggly messages, when a panic
	// was caused on purpose. So if a Pnaic was caused by a utils.ControlledPanic,
	// We'll print a nice looking error instead
	//
	// the go_throwable (ex) package helps to achieve this easy
	ex.Try(func() any {
		cmds.LetsGo()
		return nil
	}, ex.TryOpts{
		SkipWarnings: true,
		Catch: func(panicCause any) any {

			if ce, ok := interface{}(panicCause).(utils.ControlledPanic); ok {
				fmt.Println(ce.Msg)
				hadAPanic = true
				panicExitCode = ce.ExitCode
			} else {
				fmt.Println("PANIC !!!! *** Runs in circles ***")
				panic(panicCause)

			}
			return nil
		},
		Finally: func() {
			database.DeInitDB()
			if hadAPanic {
				os.Exit(panicExitCode)
			}
		},
	})
}
