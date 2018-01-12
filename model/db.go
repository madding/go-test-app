package model

import (
	"fmt"
	"log"

	"apla_test_work/config"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres" // postgrea adapter
)

var (
	// DBConn - connection to database
	DBConn *gorm.DB
)

// GormInit - init connection to database
func GormInit() error {
	conf := &config.DBConfig{}
	err := conf.Read()
	if err != nil {
		log.Fatal("Configuration cannot read, check config.ini")
	}

	DBConn, err = gorm.Open("postgres",
		fmt.Sprintf("host=localhost user=%s dbname=%s sslmode=disable password=%s",
			conf.DBUser, conf.DBName, conf.DBPass))
	if err != nil {
		return err
	}

	DBConn.AutoMigrate(&User{})

	return nil
}

// GormClose - close connection
func GormClose() error {
	if DBConn != nil {
		return DBConn.Close()
	}
	return nil
}
