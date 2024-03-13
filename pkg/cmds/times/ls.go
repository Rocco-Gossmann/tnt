package times

import (
	"fmt"
	"strings"

	"github.com/rocco-gossmann/tnt/pkg/database"
	"github.com/rocco-gossmann/tnt/pkg/utils"
	"github.com/spf13/cobra"
)

var LSCMD = cobra.Command{
	Use:   "ls",
	Short: "show all times",
	Long:  "Shows a list of all registerd times, with the latest registerd time being shown first",
	Run: func(cmd *cobra.Command, args []string) {

		times, err := database.GetTimes()
		utils.Err(err)

		sTimes := make([]string, len(times))
		for k, t := range times {
			sTimes[k] = t.String()
		}

		fmt.Println(strings.Join(sTimes, "\r\n"))

	},
}
