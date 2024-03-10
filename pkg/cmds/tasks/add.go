package tasks

import (
	"fmt"
	"strings"

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

		taskKey := database.GenerateTaskKey(taskName)
		_, err := database.ExecStatement("INSERT INTO tasks(textkey, name) VALUES (?, ?)", taskKey, taskName)

		if err != nil {
			errStr := fmt.Sprintf("%s", err)
			if strings.HasPrefix(errStr, "UNIQUE") {
				utils.Failf("task '%s' already added", taskName)
			} else {
				panic(err)
			}
		}

		fmt.Printf("added task '%s' \n", taskName)

	},
}
