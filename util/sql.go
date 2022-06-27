package util

import (
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

func GetDB() *gorm.DB {
	return db
}

//connect to database
func InitDB() (err error) {
	sqlStr := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", GetConfig().Database.DBUser, GetConfig().Database.Password, GetConfig().Database.DBHost, GetConfig().Database.DBPort, GetConfig().Database.DbName)

	db, err = gorm.Open(mysql.New(mysql.Config{
		DSN:                       sqlStr, // DSN data source name
		DefaultStringSize:         256,    // default string size
		DisableDatetimePrecision:  true,   // disable datetime Precision
		DontSupportRenameIndex:    true,   //
		DontSupportRenameColumn:   true,   //
		SkipInitializeWithVersion: false,  //
	}), &gorm.Config{})
	if err != nil {
		return err
	}
	db.Logger.LogMode(3)
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	// SetMaxIdleConns
	sqlDB.SetMaxIdleConns(10)
	// SetMaxOpenConns
	sqlDB.SetMaxOpenConns(100)
	// SetConnMaxLifetime
	sqlDB.SetConnMaxLifetime(time.Hour)
	return nil
}
