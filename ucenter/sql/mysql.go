package sql

import (
	"common/database"
	"gorm.io/gorm"
	"ucenter/config"
)

var db *gorm.DB

func InitMysql() {
	var err error
	db, err = database.ConnectMysql(config.GetConf().Mysql.DataSource)
	if err != nil {
		return
	}
}

func GetMysql() *gorm.DB {
	return db
}
