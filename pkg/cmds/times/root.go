package times

import (
	"fmt"

	"github.com/rocco-gossmann/tnt/pkg/database"
	"github.com/spf13/cobra"
)

func getTaskIdFromFlags(cmd *cobra.Command) (taskId uint) {
	taskName := cmd.Flag("task").Value.String()

	if taskName != "" {
		taskId = database.GetTaskIDByName(taskName)
	}

	return
}

func getTaskWhere(cmd *cobra.Command, prefix string) (taskWhere string) {

	taskId := getTaskIdFromFlags(cmd)

	if taskId > 0 {
		taskWhere = fmt.Sprintf("%s taskId=%d", prefix, taskId)
	}

	return
}

var TimesCMD cobra.Command = cobra.Command{
	Use:   "times {sum|ls} [-t taskName]",
	Short: "Handles everything related to Times",
	Long:  "Use this to get an Indea how much time was spend on what tasks",
}

func init() {
	TimesCMD.AddCommand(&SumCMD, &LSCMD)
	TimesCMD.PersistentFlags().StringP("task", "t", "", "filter list by task")
	TimesCMD.RegisterFlagCompletionFunc("task", database.AutoCompleteTaskList)
}
