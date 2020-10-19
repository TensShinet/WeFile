package conn

import (
	"github.com/TensShinet/WeFile/service/db/conf"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"sync"
)

var (
	dbOnce sync.Once
	db     *gorm.DB
)

func dbInit() {
	dbOnce.Do(func() {
		config := conf.GetConfig()
		var err error
		// TODO: 通过 enabled 字段支持多种数据库
		mysqlInfo := config.DB.MySQL
		db, err = gorm.Open(mysql.Open(mysqlInfo.DSN), &gorm.Config{})
		if err != nil {
			logger.Panicf("mysql open failed, for the reason:%v", err.Error())
			return
		}
		// TODO:设置连接池
		//sqlDB, err := db.DB()
		//if err != nil {
		//	logger.Panicf("mysql connect poll open failed, for the reason:%v", err.Error())
		//	return
		//}
		//sqlDB.SetConnMaxLifetime(mysqlInfo.ConnMaxLifetime)
		//sqlDB.SetMaxIdleConns(mysqlInfo.MaxIdleConnection)
		//sqlDB.SetMaxOpenConns(mysqlInfo.MaxOpenConnection)
	})
}

func GetDB() *gorm.DB {
	return db
}
