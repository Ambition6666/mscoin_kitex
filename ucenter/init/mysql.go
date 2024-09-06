package init

import (
	"common/database"
	"gorm.io/gorm"
	"ucenter/config"
)

var db *gorm.DB

func initMysql() {
	var err error
	db, err = database.ConnectMysql(config.GetConf().Mysql.DataSource)
	if err != nil {
		return
	}
}

func GetMysql() *gorm.DB {
	return db
}
