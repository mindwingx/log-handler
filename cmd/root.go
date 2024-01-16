package cmd

import (
	"context"
	"fmt"
	"github.com/mindwingx/log-handler/database/mysql"
	"github.com/spf13/cobra"
	"log"
)

var rootCmd = &cobra.Command{
	Short: "Logger main command",
	Long:  `The Logger command to manage log files.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Hand Shake!")
	},
}

func Execute(sql mysql.SqlAbstraction) {
	rootCmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		ctx := context.WithValue(cmd.Context(), "database", sql)
		cmd.SetContext(ctx)
	}

	rootCmd.PersistentPostRun = func(cmd *cobra.Command, args []string) {
		// Close the database when the command is done
		sql.Close()
	}

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}

}

func init() {
	rootCmd.AddCommand(logFileSeederCmd)
	rootCmd.AddCommand(logManagerCmd)
}
