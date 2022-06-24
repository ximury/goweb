package model

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"logger"
	"main/config"
)

func MysqlConn() *gorm.DB {
	configBase, err := config.GetChannelConfig()
	if err != nil {
		logger.Fatalf("Get config failed! err: #%v", err)
		return nil
	}
	username := configBase.Mysqlnd.Username
	password := configBase.Mysqlnd.Password
	host := configBase.Mysqlnd.Host
	port := configBase.Mysqlnd.Port
	dbname := configBase.Mysqlnd.Database
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", username, password, host, port, dbname)
	// 连接 mysql
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Fatalf("MySQL connect failed! err: #%v", err)
		return nil
	}
	// 设置数据库连接池参数
	sqlDB, _ := db.DB()
	// 设置数据库连接池最大连接数
	sqlDB.SetMaxOpenConns(100)
	// 连接池最大允许的空闲连接数，如果没有sql任务需要执行的连接数大于20，超出的连接会被连接池关闭
	sqlDB.SetMaxIdleConns(20)
	return db
}

var Db *gorm.DB

func init() {
	Db = MysqlConn()
}
