package cmd

import (
	"fmt"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/fatih/color"
	src "github.com/mindwingx/log-handler"
	"github.com/mindwingx/log-handler/utils"
	"github.com/spf13/cobra"
	"log"
	"math/rand"
	"os"
	"sync"
	"time"
)

var logFileSeederCmd = &cobra.Command{
	Use:   "logs:seed",
	Short: "generate dummy log files",
	Run: func(cmd *cobra.Command, args []string) {
		evaluateLogsDirExistence()

		logger := initLogGenerator(cmd.Context().Value("log_counter").(int))
		logger.wg.Add(logger.wgCap)

		for i := 1; i <= logger.logCounter; i++ {
			// pass the required variables as a batch by a closure function
			go logger.dummyLogFileGenerator()
		}

		logger.wg.Wait()

		color.Cyan(fmt.Sprintf(
			"seeding of %d was done within %d miliseconds\n",
			logger.logCounter, time.Since(logger.start).Milliseconds(),
		))
	},
}

//HELPER METHODS

func evaluateLogsDirExistence() {
	if err := os.MkdirAll(fmt.Sprintf("%s/logs", src.Root()), os.ModePerm); err != nil {
		log.Fatal("Error creating logs directory:", err)
		// todo: it is recommended to trace error by the Sentry
	}
}

// LOGGER SCOPE

type logGenerator struct {
	logCounter int
	wgCap      int
	wg         sync.WaitGroup
	mx         sync.Mutex
	start      time.Time
}

func initLogGenerator(logCounter int) *logGenerator {
	return &logGenerator{
		logCounter: logCounter,
		wgCap:      logCounter,
		wg:         sync.WaitGroup{},
		mx:         sync.Mutex{},
		start:      time.Now(),
	}
}

func (lg *logGenerator) dummyLogFileGenerator() {
	defer lg.wg.Done()

	//primary variables
	var dummyLog string

	lg.mx.Lock()
	filePath := fmt.Sprintf("%s/logs/%d.log", src.Root(), time.Now().UnixNano())
	//inquire file existence
	if _, statErr := os.Stat(filePath); statErr == nil {
		// if file exists, decrease the counter to create new file
		color.Yellow("log file exists:\n", filePath, "\n")
		lg.logCounter++
		lg.mx.Unlock()
		return
	}
	lg.mx.Unlock()

	//init log file
	file, err := os.Create(filePath)
	if err != nil {
		// just a snapshot in STDOUT
		// todo: it is recommended to trace error by the Sentry
		fmt.Println(err.Error())
	}

	defer file.Close()

	for i := 0; i < rand.Intn(6); i++ {
		//generate dummy log
		randLogRange := rand.Intn(5)
		dummyLog += fmt.Sprintf(
			"[%s] - %s: %s\n",
			time.Now().Format(utils.TimestampLayout),
			utils.LogLevels[randLogRange],
			gofakeit.HackerPhrase(),
		)
	}

	// append multiple log along with current log file by inserting new ones
	// to write is done outside the loop to avoid extra probable read and write
	_, err = file.WriteString(dummyLog)
	if err != nil {
		// just a snapshot in STDOUT
		// todo: it is recommended to trace error by the Sentry
		fmt.Println(err.Error())
	}
}
