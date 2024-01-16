package cmd

import (
	"fmt"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/fatih/color"
	sdstudio "github.com/mindwingx/log-handler"
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
		// seeding primary variables
		wgCap, logFilesCounter := 1000, 1000
		st := time.Now()

		// create the logs directory if not exists
		if err := os.MkdirAll(fmt.Sprintf("%s/logs", sdstudio.Root()), os.ModePerm); err != nil {
			log.Fatal("Error creating logs directory:", err)
			// todo: it is recommended to trace error by the Sentry
		}

		//the sync package variables
		wg := sync.WaitGroup{}
		mx := sync.Mutex{}
		wg.Add(wgCap)

		for i := 1; i <= logFilesCounter; i++ {
			// pass the required variables as a batch by a closure function
			go func(wg *sync.WaitGroup, mx *sync.Mutex, logFileCounter *int) {
				dummyLogFileGenerator(wg, mx, &logFilesCounter)
			}(&wg, &mx, &logFilesCounter)
		}

		wg.Wait()

		fmt.Printf("seeding done within %d miliseconds\n", time.Since(st).Milliseconds())
	},
}

//HELPER METHODS

func dummyLogFileGenerator(wg *sync.WaitGroup, mx *sync.Mutex, logFileCounter *int) {
	defer wg.Done()

	//primary variables
	var dummyLog string

	mx.Lock()
	filePath := fmt.Sprintf("%s/logs/%d.log", sdstudio.Root(), time.Now().UnixNano())
	//inquire file existence
	if _, statErr := os.Stat(filePath); statErr == nil {
		// if file exists, decrease the counter to create new file
		color.Yellow("log file exists:\n", filePath, "\n")
		*logFileCounter++
		mx.Unlock()
		return
	}
	mx.Unlock()

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
