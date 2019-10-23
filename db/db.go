package db

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
	"sync"
)

type db struct {
	Mysql *gorm.DB
}

var dbInstance *db
var dbInstanceOnce sync.Once

func InitDB(mySqlInfo string) {
	dbInstanceOnce.Do(func() {
		mysqlDB, err := gorm.Open("mysql", mySqlInfo)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"mySqlInfo": mySqlInfo,
				"err": err.Error(),
			}).Fatal("Mysql init failed.")
		}
		dbInstance = &db{
			Mysql: mysqlDB,
		}
	})
}

func GetDB() *db {
	return dbInstance
}

func (instance *db) Close() {
	instance.Mysql.Close()
}
