package main

import (
	"fmt"
	"github.com/alexflint/go-filemutex"
	"github.com/mindwingx/log-handler/cmd"
	"github.com/mindwingx/log-handler/database/mysql"
	"github.com/mindwingx/log-handler/registry"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	mutex, err := filemutex.New("/tmp/log-handler.lock")
	if err != nil {
		log.Fatal("/tmp directory does not exist or lock file cannot be created!")
	}

	errLock := mutex.TryLock()
	if errLock != nil {
		log.Fatal("This program is already running on this server!")
		return
	}

	mutex.Lock()

	reg := registry.NewViper()
	reg.InitRegistry()

	db := mysql.NewSql(reg)
	db.InitSql()

	cmd.Execute(db)

	fmt.Println("Do interrupt!")
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
	fmt.Println("Interrupt is detected")

	mutex.Unlock()
}
