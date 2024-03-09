package cmds

import (
	"fmt"

	"github.com/rocco-gossmann/tnt/pkg/database"
	"github.com/spf13/cobra"
)

var StopCmd = cobra.Command{
	Use:                   "stop",
	Short:                 "Stops all timers",
	Long:                  "Stops all timers",
	DisableFlagsInUseLine: true,
	Args:                  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		database.FinishCurrentlyRunningTimes()
		fmt.Println("done")
	},
}
