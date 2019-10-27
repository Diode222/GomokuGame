package db

import (
	"GomokuGame/app/conf"
	"github.com/go-redis/redis"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
	"sync"
	"time"
)

type DB struct {
	Lock  sync.Locker
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
				"err":       err.Error(),
			}).Fatal("Mysql init failed.")
		}
		redisCli := redis.NewClient(&redis.Options{
			Addr:        conf.REDIS_ADDR,
			Password:    conf.REDIS_PASSWORD,
			DB:          conf.REDIS_DB_NAME,
			DialTimeout: 5 * time.Second,
		})
		dbInstance = &DB{
			Lock:  &sync.Mutex{},
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
