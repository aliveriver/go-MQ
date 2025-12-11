package common

import (
	"fmt"

	"go-MQ/entity"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

func getModels() []interface{} {
	return []interface{}{
		&entity.User{},
	}
}

func InitDB() *gorm.DB {
	if db != nil {
		logrus.Fatal("DB is already initialized")
	}

	host := viper.GetString("database.host")
	port := viper.GetString("database.port")
	database := viper.GetString("database.database")
	username := viper.GetString("database.username")
	password := viper.GetString("database.password")
	charset := viper.GetString("database.charset")
	loc := viper.GetString("database.loc")
	args := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=true&loc=%s",
		username,
		password,
		host,
		port,
		database,
		charset,
		loc,
	)

	tempdb, err := gorm.Open(mysql.Open(args), &gorm.Config{})
	if err != nil {
		logrus.Fatal("failed to connect database,err:" + err.Error() + ",args:" + args)
	}

	db = tempdb

	if err := db.AutoMigrate(getModels()...); err != nil {
		logrus.Fatal("Failed to migrate database:", err)
	}
	return db
}

func GetDB() *gorm.DB {
	if db == nil {
		panic("DB is nil")
	}
	return db
}
