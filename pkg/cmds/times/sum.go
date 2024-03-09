package times

import (
	"fmt"

	"github.com/rocco-gossmann/tnt/pkg/database"
	"github.com/rocco-gossmann/tnt/pkg/utils"
	"github.com/spf13/cobra"
)

var SumCMD = cobra.Command{
	Use:   "sum",
	Short: "shows how much time was taken on what task",
	Run: func(cmd *cobra.Command, args []string) {

		res, err := database.QueryStatement(`
			SELECT 
				ta.name, 
				time(sum(unixepoch(ti.end) - unixepoch(ti.start)), "unixepoch") total
			FROM times ti 
				LEFT JOIN tasks ta ON ti.taskId = ta.id
			WHERE end IS NOT NULL GROUP by taskId
		`)

		utils.Err(err)

		for res.Next() {
			var name, total string
			err = res.Scan(&name, &total)
			utils.Err(err)
			fmt.Printf(" %s | %s Hours \n", name, total)
		}

	},
}
