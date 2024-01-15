package main

import (
	"github.com/alexflint/go-filemutex"
	"github.com/mindwingx/log-handler/cmd"
	"github.com/mindwingx/log-handler/database/mysql"
	"github.com/mindwingx/log-handler/registry"
	"github.com/mindwingx/log-handler/utils"
	"log"
)

func main() {
	mutex, err := filemutex.New(utils.TmpLockFile)

	if err != nil {
		log.Fatal("/tmp directory does not exist or lock file cannot be created!")
	}

	errLock := mutex.TryLock()
	if errLock != nil {
		log.Fatal("This program is already running on this server!")
		return
	}

	err = mutex.Lock()
	if err != nil {
		log.Fatal(err)
	}

	reg := registry.NewViper()
	reg.InitRegistry()

	db := mysql.NewSql(reg)
	db.InitSql()
	db.Migrate()

	cmd.Execute(reg, db)
	err = mutex.Unlock()
	if err != nil {
		log.Fatal(err)
	}
}
