package mysql

import (
	"fmt"
	src "github.com/mindwingx/log-handler"
	"github.com/mindwingx/log-handler/registry"
	"github.com/mindwingx/log-handler/utils"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"sort"
	"time"
)

type SqlAbstraction interface {
	InitSql()
	Migrate()
	Sql() *Sql
	Close()
	//Seed()
}

type (
	Sql struct {
		config config
		DB     *gorm.DB
	}

	config struct {
		Debug              bool   `mapstructure:"DEBUG"`
		Host               string `mapstructure:"HOST"`
		Port               string `mapstructure:"PORT"`
		Username           string `mapstructure:"USERNAME"`
		Password           string `mapstructure:"PASSWORD"`
		RootPassword       string `mapstructure:"ROOT_PASSWORD"`
		Database           string `mapstructure:"DATABASE"`
		Ssl                string `mapstructure:"SSL"`
		MaxIdleConnections int    `mapstructure:"MAXIDLECONNECTIONS"`
		MaxOpenConnections int    `mapstructure:"MAXOPENCONNECTIONS"`
		MaxLifetimeSeconds int    `mapstructure:"MAXLIFETIMESECONDS"`
	}
)

func NewSql(registry registry.RegAbstraction) SqlAbstraction {
	database := new(Sql)
	registry.Parse(&database.config)
	return database
}

func (g *Sql) InitSql() {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=true",
		g.config.Username,
		g.config.Password,
		g.config.Host,
		g.config.Port,
		g.config.Database,
	)
	database, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		SkipDefaultTransaction: true,
		Logger:                 g.newGormLog(),
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	})
	if err != nil {

		log.Fatal(err)
	}

	sqlDatabase, err := database.DB()
	if err != nil {
		log.Fatal(err)
	}

	if g.config.MaxIdleConnections != 0 {
		sqlDatabase.SetMaxIdleConns(g.config.MaxIdleConnections)
	}

	if g.config.MaxOpenConnections != 0 {
		sqlDatabase.SetMaxOpenConns(g.config.MaxOpenConnections)
	}

	if g.config.MaxLifetimeSeconds != 0 {
		sqlDatabase.SetConnMaxLifetime(time.Second * time.Duration(g.config.MaxLifetimeSeconds))
	}

	if g.config.Debug {
		database = database.Debug()
	}

	g.DB = database
}

func (g *Sql) Migrate() {
	path := fmt.Sprintf("%s/database/mysql/migrations", src.Root())

	// Open the directory
	dir, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
		return
	}

	defer dir.Close()

	// Read the directory contents
	fileInfos, err := dir.Readdir(-1)
	if err != nil {
		log.Fatal(err)
		return
	}

	// Sort the entries alphabetically by name - Sql file order by numeric(01, 02, etc)
	sort.Slice(fileInfos, func(i, j int) bool {
		return fileInfos[i].Name() < fileInfos[j].Name()
	})

	// Iterate over the file info slice and print the file names
	for _, fileInfo := range fileInfos {
		if fileInfo.Mode().IsRegular() {
			if err = g.DB.Exec(g.parseSqlFile(path, fileInfo)).Error; err != nil {
				log.Fatal(err)
			}
		}
	}
}

func (g *Sql) Sql() *Sql {
	return g
}

func (g *Sql) Close() {
	sql, err := g.DB.DB()
	if err != nil {
		log.Fatal(err)
		return
	}

	sql.Close()
}

// HELPER METHOD

func (g *Sql) newGormLog() logger.Interface {
	return logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             utils.SlowSqlThreshold * time.Second, // Slow SQL threshold
			LogLevel:                  logger.Warn,                          // Log level
			IgnoreRecordNotFoundError: false,                                // Ignore ErrRecordNotFound error for logger
			Colorful:                  true,                                 // Disable color
		})
}

func (g *Sql) parseSqlFile(path string, fileInfo os.FileInfo) string {
	sqlFile := fmt.Sprintf("%s/%s", path, fileInfo.Name())
	sqlBytes, err := os.ReadFile(sqlFile)
	if err != nil {
		log.Fatal(err)
	}
	// Convert SQL file contents to string
	return string(sqlBytes)
}
