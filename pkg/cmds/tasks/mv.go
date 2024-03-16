package tasks

import (
	"log"

	"github.com/rocco-gossmann/tnt/pkg/database"
	"github.com/rocco-gossmann/tnt/pkg/utils"
	"github.com/spf13/cobra"
)

var MVCMD cobra.Command = cobra.Command{
	Use:               "mv taskName newTaskName",
	Aliases:           []string{"rename"},
	Short:             "raname a Task",
	Long:              "Rename a Task",
	Args:              cobra.ExactArgs(2),
	ValidArgsFunction: database.AutoCompleteTaskList,

	Run: func(cmd *cobra.Command, args []string) {
		//TODO: Implement

		var targetName, newName string = args[0], args[1]
		taskId := database.GetTaskIDByName(targetName)

		res, err := database.RenameTask(taskId, newName)

		log.Println("Result:", res)
		log.Println("Error:", err)

		utils.Err(err)

	},
}
