package cmd

import (
	"bufio"
	"fmt"
	"github.com/da0x/golang/olog"
	"github.com/fatih/color"
	src "github.com/mindwingx/log-handler"
	"github.com/mindwingx/log-handler/constants"
	"github.com/mindwingx/log-handler/database/models"
	"github.com/mindwingx/log-handler/database/mysql"
	"github.com/mindwingx/log-handler/utils"
	"github.com/spf13/cobra"
	"gorm.io/gorm"
	"io/fs"
	"log"
	"os"
	"regexp"
	"strings"
	"time"
)

type logDetailSnapshot struct {
	LogID    uint64    `json:"log_id"`
	LogName  string    `json:"log_name"`
	LogLevel string    `json:"log_level"`
	LogDate  time.Time `json:"log_date"`
}

var logManagerCmd = &cobra.Command{
	Use:   "logs:run",
	Short: "trigger log handler",
	Run: func(cmd *cobra.Command, args []string) {
		logs, scheduler, logLevelCriteria, db := primaryVariables(cmd)

		for _, log := range logs {
			file := fmt.Sprintf("%s/logs/%s", src.Root(), log.Name())
			if _, statErr := os.Stat(file); statErr != nil {
				color.Red("\ninvalid file: ", log.Name())
				// todo: it is recommended to trace error by the Sentry
				break
			}

			reader, err := parsedLogReader(file)
			if err != nil {
				color.Red("\nfailed to parse file: ", log.Name())
				// todo: it is recommended to trace error by the Sentry
				break
			}

			var task utils.Task = func() {
				// prepend primary log model
				plm := models.PrimaryLogs{Uuid: utils.NewUuid(), File: log.Name(), ProcessState: constants.Init}

				if newLogErr := db.Create(&plm).Error; newLogErr != nil {
					color.Red("log details db insert failure:", err.Error())
					// todo: it is recommended to trace error by the Sentry
				}

				if len(reader) > 0 {
					prependLogDetails(reader, &plm, logLevelCriteria)

					if errDb := db.Omit("updated_at", "deleted_at").Save(&plm).Error; errDb != nil {
						color.Red("log details db insert failure:", err.Error())
						// todo: it is recommended to trace error by the Sentry
					} else {
						if len(plm.Details) > 0 {
							go printLogDetailsInStd(log.Name(), &plm)
						}
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

func primaryVariables(cmd *cobra.Command) (logs []os.DirEntry, scheduler *utils.Scheduler, logCriteria string, db *gorm.DB) {
	logs = logsDirScanner()
	scheduler = utils.InitScheduler(cmd.Context().Value("worker_count").(int))
	logCriteria = cmd.Context().Value("log_level_criteria").(string)
	db = cmd.Context().Value("database").(mysql.SqlAbstraction).Sql().DB
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

func prependLogDetails(reader [][]string, primaryLog *models.PrimaryLogs, logLevelCriteria string) {
	for _, trace := range reader {
		date, level, commit := extractLogInfo(trace[0])

		// handle log level criteria before transformation
		if logLevelCriteria != "*" && logLevelCriteria != level {
			break
		}

		timestamp, _ := time.Parse(constants.TimestampLayout, date)

		detail := models.LogDetails{
			CreatedAt: timestamp,
			LogLevel:  level,
			Log:       commit,
		}

		primaryLog.Details = append(primaryLog.Details, detail)
	}

	primaryLog.ProcessState = constants.Done
}

func printLogDetailsInStd(logName string, primaryLog *models.PrimaryLogs) {
	var logsSnapshot []logDetailSnapshot

	for _, detail := range primaryLog.Details {
		logsSnapshot = append(logsSnapshot, logDetailSnapshot{
			LogID:    primaryLog.ID,
			LogName:  logName,
			LogLevel: detail.LogLevel,
			LogDate:  detail.CreatedAt,
		})
	}

	olog.Print(logsSnapshot)
}
