package tasks

import (
	"github.com/spf13/cobra"
)

var TaskCMD cobra.Command = cobra.Command{
	Use:   "tasks {ls|add|rm}",
	Short: "Handles everything related to Tasks",
	Long:  "Use this to create new Tasks to track or remove tasks you don't want to track anymore",
}

func init() {
	TaskCMD.AddCommand(&AddCMD, &LSCMD, &RMCMD, &MVCMD)
}
