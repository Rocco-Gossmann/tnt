package tasks

import (
	"fmt"

	"github.com/rocco-gossmann/tnt/pkg/database"
	"github.com/rocco-gossmann/tnt/pkg/utils"
	"github.com/spf13/cobra"
)

var LSCMD = cobra.Command{
	Use:                   "ls",
	Short:                 "list",
	Long:                  "Lists all previously entered tasks.",
	DisableFlagsInUseLine: true,
	Run: func(cmd *cobra.Command, args []string) {
		tasks, err := database.GetTaskList()
		utils.Err(err)

		for _, task := range tasks {
			fmt.Println(task.Name)
		}
	},
}
