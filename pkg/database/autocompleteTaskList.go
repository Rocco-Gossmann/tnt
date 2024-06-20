package database

import (
	"github.com/rocco-gossmann/tnt/pkg/utils"
	"github.com/spf13/cobra"
)

func AutoCompleteTaskList(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {

	tasks, err := GetTaskList("")
	utils.Err(err)

	lst := TaskList(tasks).ExtractTaskListNames()
	return lst, cobra.ShellCompDirectiveNoFileComp

}
