package pgsql

import (
	"database/sql"
	"errors"
	"log"
	"os"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/plugin/opentelemetry/tracing"
)

var (
	sqlDb *sql.DB
	db    *gorm.DB
)

var serverLogger = logger.New(
	log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
	logger.Config{
		LogLevel:                  logger.Warn, // Log level
		IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
		ParameterizedQueries:      true,        // Don't include params in the SQL log
		Colorful:                  false,       // Disable color
	},
)

func Init(cfg *Config) (err error) {
	dbInstance, err := gorm.Open(
		postgres.New(postgres.Config{
			DSN: cfg.dsn(),
			// disables implicit prepared statement usage
			PreferSimpleProtocol: true,
		}),
		&gorm.Config{
			Logger: serverLogger,
		},
	)
	if err != nil {
		return
	}

	if err := dbInstance.Use(tracing.NewPlugin()); err != nil {
		panic(err)
	}

	if sqlDb, err = dbInstance.DB(); err != nil {
		return
	}

	sqlDb.SetConnMaxLifetime(cfg.getLifeDuration())
	sqlDb.SetMaxOpenConns(cfg.MaxConn)
	sqlDb.SetMaxIdleConns(cfg.MaxIdle)

	db = dbInstance

	return
}

func DB() *gorm.DB {
	return db
}

// AutoMigrate automatically migrates the given table.
func AutoMigrate(table interface{}) error {
	// pointer non-nil check
	if db == nil {
		err := errors.New("empty pointer of *gorm.DB")
		logrus.Error(err.Error())

		return err
	}

	return db.AutoMigrate(table)
}
