package times

import (
	"database/sql"
	"fmt"

	"github.com/rocco-gossmann/tnt/pkg/database"
	"github.com/rocco-gossmann/tnt/pkg/utils"
	"github.com/spf13/cobra"
)

var LSCMD = cobra.Command{
	Use:   "ls",
	Short: "show all times",
	Long:  "Shows a list of all registerd times, with the latest registerd time being shown first",
	Run: func(cmd *cobra.Command, args []string) {

		res, err := database.QueryStatement(`
			SELECT 
				ta.name, 
				ti.start, 
				ti.end, 
				time(unixepoch(ti.end) - unixepoch(ti.start), "unixepoch") duration
			FROM times ti
			LEFT JOIN tasks ta ON ti.taskId = ta.id
			ORDER BY start DESC;
		`)

		utils.Err(err)

		for res.Next() {
			var name, total, start, end sql.NullString
			err = res.Scan(&name, &start, &end, &total)
			utils.Err(err)

			if !end.Valid {
				end.String = "* running *"
			}

			if total.Valid {
				total.String += " Hours"
			}

			fmt.Printf(" %s | %s | %s | %s \n", name.String, start.String, end.String, total.String)
		}

	},
}
