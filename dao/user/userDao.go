package user

import (
	"GomokuGame/app/conf"
	"GomokuGame/db"
	"GomokuGame/model"
	"GomokuGame/utils/jwt"
	"context"
	"encoding/json"
	"errors"
	"strings"
)

type UserDao struct {
	DBInstance *db.DB
}

func NewUserDao(dbInstance *db.DB) *UserDao {
	return &UserDao{
		DBInstance: dbInstance,
	}
}

func (u *UserDao) Login(ctx context.Context, userName string, password string) (string, error) {
	userInfoBinary, err := u.DBInstance.Redis.Get(userName).Result()
	if err == nil && userInfoBinary != "" {
		u.setUserTokenPairInRedis(userInfoBinary)
		return userInfoBinary, nil
	}

	userItem := &model.UserItem{}
	u.DBInstance.Mysql.Table(conf.USER_TABLE_NAME).Where("user_name = ?", userName).Where("password = ?", password).Find(userItem)
	if userItem.UserName == "" || userItem.Password == "" {
		return "", errors.New("wrong user info")
	}

	bytes, err := json.Marshal(userItem)
	if err != nil {
		return "", errors.New("UserItem marshal failed")
	}

	userInfoBinary = string(bytes)
	token := u.setUserTokenPairInRedis(userInfoBinary)
	return token, nil
}

func (u *UserDao) Register(ctx context.Context, userName string, password string, warehouseAddr string) error {
	if !u.availableWarehouseAddr(warehouseAddr) {
		return errors.New("Git warehouse address is not right.")
	}

	users := []*model.UserItem{}
	u.DBInstance.Mysql.Table(conf.USER_TABLE_NAME).Where("user_name = ?", userName).Find(&users)
	if len(users) > 0 {
		return errors.New("register failed")
	}

	times := 0
	for times < 10 {
		err := db.GetDB().Mysql.Table(conf.USER_TABLE_NAME).Create(&model.UserItem{
			UserName:      userName,
			Password:      password,
			WarehouseAddr: warehouseAddr,
		}).Error
		if err == nil {
			break
		}
		times++
	}
	if times == 10 {
		return errors.New("register failed")
	}

	return nil
}

func (u *UserDao) availableWarehouseAddr(addr string) bool {
	if !strings.HasPrefix(addr, "git@github.com:") || !strings.HasSuffix(addr, "/GomokuGameImpl.git") {
		return false
	}
	return true
}

func (u *UserDao) setUserTokenPairInRedis(userInfo string) string {
	token := jwt.GetToken()
	u.DBInstance.Redis.Set(token, userInfo, conf.USER_INFO_TTL)
	return token
}
