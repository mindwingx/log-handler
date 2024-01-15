package cmd

import (
	"context"
	"fmt"
	"github.com/mindwingx/log-handler/database/mysql"
	"github.com/mindwingx/log-handler/registry"
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

type config struct {
	LogCounter  int `mapstructure:"SEEDER_LOG_COUNTER"`
	WorkerCount int `mapstructure:"RUNNER_WORKER_COUNT"`
}

func Execute(registry registry.RegAbstraction, sql mysql.SqlAbstraction) {
	var conf config
	registry.Parse(&conf)

	rootCmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		cmd.SetContext(context.WithValue(cmd.Context(), "database", sql))
		cmd.SetContext(context.WithValue(cmd.Context(), "log_counter", conf.LogCounter))
		cmd.SetContext(context.WithValue(cmd.Context(), "worker_count", conf.WorkerCount))
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
