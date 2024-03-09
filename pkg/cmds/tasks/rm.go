package tasks

import (
	"fmt"

	"github.com/rocco-gossmann/tnt/pkg/database"
	"github.com/rocco-gossmann/tnt/pkg/utils"
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

		_, err := database.ExecStatement("DELETE FROM times WHERE taskId = ?", taskId)
		utils.Err(err)

		_, err = database.ExecStatement("DELETE FROM tasks WHERE id = ?", taskId)
		utils.Err(err)

		fmt.Printf("removed task '%s' \n", args[0])
	},
}
