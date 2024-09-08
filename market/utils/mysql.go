package utils

import (
	"common/database"
	"gorm.io/gorm"
	"market/config"
)

var db *gorm.DB

func initMysql() {
	var err error
	db, err = database.ConnectMysql(config.GetConf().Mysql.DataSource)
	if err != nil {
		panic(err)
	}
}

func GetMysql() *gorm.DB {
	return db
}
