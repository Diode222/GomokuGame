package db

import (
	"github.com/go-redis/redis"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
	"sync"
	"time"
)

type DB struct {
	Mysql *gorm.DB
	Redis *redis.Client
}

var dbInstance *DB
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
		redisCli := redis.NewClient(&redis.Options{
			Addr:               "139.155.46.62:6379",
			Password:           "",
			DB:                 0,
			DialTimeout: 5 * time.Second,
		})
		dbInstance = &DB{
			Mysql: mysqlDB,
			Redis: redisCli,
		}
	})
}

func GetDB() *DB {
	return dbInstance
}

func (instance *DB) Close() {
	instance.Mysql.Close()
	instance.Redis.Close()
}
