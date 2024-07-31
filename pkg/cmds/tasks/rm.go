package tasks

import (
	"fmt"

	"github.com/rocco-gossmann/tnt/pkg/database"
	"github.com/spf13/cobra"
)

var RMCMD cobra.Command = cobra.Command{
	Use:                   "rm \"Taskname\"",
	Short:                 "remove",
	Long:                  "Removes a Task from the List of previously entered tasks.",
	DisableFlagsInUseLine: true,
	Args:                  cobra.ExactArgs(1),
	ValidArgsFunction:     database.AutoCompleteTaskList,
	Run: func(cmd *cobra.Command, args []string) {

		taskId := database.GetTaskIDByName(args[0])
		database.DropTask(taskId)

		fmt.Printf("removed task '%s' \n", args[0])
	},
}
