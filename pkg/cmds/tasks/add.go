package tasks

import (
	"fmt"

	"github.com/rocco-gossmann/tnt/pkg/database"
	"github.com/rocco-gossmann/tnt/pkg/utils"
	"github.com/spf13/cobra"
)

var AddCMD cobra.Command = cobra.Command{
	Use:                   "add \"Taskname\"",
	Short:                 "Adds a new Task",
	Args:                  cobra.ExactArgs(1),
	DisableFlagsInUseLine: true,
	Run: func(cmd *cobra.Command, args []string) {

		taskName := args[0]
		err := database.AddTask(args[0])

		if err != nil {
			if database.IsUniqueContraintError(err) {
				utils.Failf("task '%s' already added", taskName)

			} else {
				panic(err)

			}
		}

		fmt.Printf("added task '%s' \n", taskName)

	},
}
