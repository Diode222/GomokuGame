package main

import (
	"GomokuGame/app/conf"
	"GomokuGame/db"
	"GomokuGame/kube"
	"GomokuGame/model"
	"GomokuGame/router"
	"fmt"
)

func main() {
	sqlInfo := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=true", conf.MYSQL_USERNAME, conf.MYSQL_PASSWORD, conf.MYSQL_IP, conf.MYSQL_PORT, conf.MYSQL_DBNAME)
	db.InitDB(sqlInfo)
	defer db.GetDB().Close()
	kube.InitClientset()
	db.GetDB().Mysql.AutoMigrate(&model.UserItem{}, &model.MatchResultItem{})
	router.InitRouter()
}
