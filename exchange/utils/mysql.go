package utils

import (
	"common/database"
	"exchange/config"
	"gorm.io/gorm"
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
