package cmds

import (
	"fmt"
	"log"

	"github.com/rocco-gossmann/tnt/pkg/database"
	"github.com/rocco-gossmann/tnt/pkg/utils"
	"github.com/spf13/cobra"
)

var SwitchCMD = cobra.Command{
	Use:                   "switch \"taskname\"",
	Short:                 "Switch to a task",
	Aliases:               []string{"s", "start"},
	Long:                  "Switches to or Starts a different task",
	DisableFlagsInUseLine: true,
	Args:                  cobra.ExactArgs(1),
	ValidArgsFunction:     database.AutoCompleteTaskList,
	Run: func(cmd *cobra.Command, args []string) {

		taskId := database.GetTaskIDByName(args[0])

		if taskId == 0 {
			utils.Failf("task '%s' not found", args[0])
		}

		if database.TimedTaskIsRunning(taskId) {
			fmt.Println("** task already running **")
			return
		}

		database.FinishCurrentlyRunningTimes()
		id := database.StartNewTime(taskId)

		log.Printf("new timer with Id: %d\n", id)
		fmt.Printf("timer started")
	},
}
