package times

import (
	"fmt"

	"github.com/rocco-gossmann/tnt/pkg/database"
	"github.com/spf13/cobra"
)

var SumCMD = cobra.Command{
	Use:   "sum",
	Short: "shows how much time was taken on what task",
	Run: func(cmd *cobra.Command, args []string) {

		sums := database.GetTimeSums(getTaskIdFromFlags(cmd))

		for _, sum := range sums {
			fmt.Printf(" %s | %s Hours \n", sum.Name, sum.Total)
		}

	},
}
