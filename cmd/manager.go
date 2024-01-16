package cmd

import (
	"fmt"
	"github.com/mindwingx/log-handler/database/models"
	"github.com/spf13/cobra"
)

var logManagerCmd = &cobra.Command{
	Use:   "logs:run",
	Short: "trigger log handler",
	Run: func(cmd *cobra.Command, args []string) {
		//dbConn := cmd.Context().Value("database").(mysql.SqlAbstraction)

		var dbs []models.PrimaryLogs
		//_ = dbConn.Sql().DB.Raw(`SELECT * FROM clients WHERE id = ?`, 1).Scan(&dbs).Error
		//fmt.Println("Executing child command...\n", dbs[0].Email)
		fmt.Println("Executing child command...\n", dbs)
	},
}
