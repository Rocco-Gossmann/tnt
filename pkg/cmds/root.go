package cmds

import (
	"fmt"
	"io"
	"log"

	"github.com/rocco-gossmann/tnt/pkg/database"
	"github.com/rocco-gossmann/tnt/pkg/env"
	"github.com/rocco-gossmann/tnt/pkg/cmds/tasks"
	"github.com/rocco-gossmann/tnt/pkg/cmds/times"
	"github.com/spf13/cobra"
)

var MyCMD = cobra.Command{
	Use: "tnt {tasks|s|start|switch|stop|times} [-v] [-h]",

	PersistentPreRun: func(cmd *cobra.Command, args []string) {

		if cmd.Flag("debug").Value.String() == "false" {
			log.SetOutput(io.Discard)
			log.Println("--debug set => enable logging")
		}

		database.InitDB("")
	},

	RunE: func(cmd *cobra.Command, args []string) error {
		if cmd.Flag("version").Value.String() == "true" {
			fmt.Println(env.Version)
			return nil
		}

		return fmt.Errorf("unknown command")
	},

	DisableFlagsInUseLine: true,
}

func init() {
	MyCMD.PersistentFlags().Bool("debug", false, "Enable Debug-Log output")
	MyCMD.PersistentFlags().BoolP("version", "v", false, "Prints the version number of Tasks n' Times")

	// Add all the Sub-Commands
	MyCMD.AddCommand(&tasks.TaskCMD, &SwitchCMD, &StopCmd, &times.TimesCMD)
}

func LetsGo() {
	MyCMD.Execute()
}
