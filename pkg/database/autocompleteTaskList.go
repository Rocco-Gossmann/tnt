package database

import (
	"github.com/rocco-gossmann/tnt/pkg/utils"
	"github.com/spf13/cobra"
)

func AutoCompleteTaskList(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {

	lst, err := GetTaskList()
	utils.Err(err)

	return lst, cobra.ShellCompDirectiveNoFileComp

}
