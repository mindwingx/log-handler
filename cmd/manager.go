package cmd

import (
	"bufio"
	"fmt"
	src "github.com/mindwingx/log-handler"
	"github.com/mindwingx/log-handler/utils"
	"github.com/spf13/cobra"
	"io/fs"
	"log"
	"os"
	"regexp"
	"strings"
	"time"
)

var logManagerCmd = &cobra.Command{
	Use:   "logs:run",
	Short: "trigger log handler",
	Run: func(cmd *cobra.Command, args []string) {
		logs := logsDirScanner()
		scheduler := utils.InitScheduler(cmd.Context().Value("worker_count").(int))
		//db := cmd.Context().Value("database").(mysql.SqlAbstraction)

		fmt.Println("\n", len(logs))
		repeat := strings.Repeat("-", 20)
		fmt.Println(repeat, repeat)

		for _, log := range logs {
			file := fmt.Sprintf("%s/logs/%s", src.Root(), log.Name())
			if _, statErr := os.Stat(file); statErr != nil {
				fmt.Println("\ninvalid file!\n", file)
				break
			}

			reader, err := parsedLogReader(file)
			if err != nil {
				break
			}

			var task utils.Task = func() {
				if len(reader) > 0 {
					for _, trace := range reader {
						date, level, commit := extractLogInfo(trace[0])
						fmt.Println(date, level, commit)
						time.Sleep(time.Second * 2)
					}
				}
			}

			scheduler.AddTask(task)
		}

		scheduler.Start()
	},
}

// HELPER METHODS

func logsDirScanner() (logs []os.DirEntry) {
	dir := fmt.Sprintf("%s/logs", src.Root())
	logs, err := fs.ReadDir(os.DirFS(dir), ".")
	if err != nil {
		log.Fatal("Error reading directory:", err)
	}
	return
}

func parsedLogReader(path string) (lines [][]string, err error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	err = scanner.Err()

	for scanner.Scan() {
		s := scanner.Text()
		lines = append(lines, strings.Split(s, "\n"))
	}

	return
}

func extractLogInfo(input string) (date string, level string, log string) {
	re := regexp.MustCompile(`\[(.*?)\] - (\w+): (.*)`)

	matches := re.FindStringSubmatch(input)
	date, level, log = matches[1], matches[2], matches[3]
	return
}
