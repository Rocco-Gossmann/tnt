package times

import (
	"github.com/spf13/cobra"
)

var TimesCMD cobra.Command = cobra.Command{
	Use:   "times {sum|ls}",
	Short: "Handles everything related to Times",
	Long:  "Use this to get an Indea how much time was spend on what tasks",
}

func init() {
	TimesCMD.AddCommand(&SumCMD, &LSCMD)
}
