package user

import (
	"GomokuGame/app/conf"
	"GomokuGame/db"
	"GomokuGame/model"
	"GomokuGame/utils/jwt"
	"context"
	"encoding/json"
	"errors"
	"github.com/sirupsen/logrus"
	"strings"
	"time"
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
	userItem := &model.UserItem{}
	u.DBInstance.Mysql.Table(conf.USER_TABLE_NAME).Where("user_name = ?", userName).Where("password = ?", password).Find(userItem)
	if userItem.UserName == "" || userItem.Password == "" {
		return "", errors.New("wrong user info")
	}

	bytes, err := json.Marshal(userItem)
	if err != nil {
		return "", errors.New("UserItem marshal failed")
	}

	userInfoBinary := string(bytes)
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
	userInfo := &model.UserItem{
		UserName:      userName,
		Password:      password,
		WarehouseAddr: warehouseAddr,
	}
	for times < 10 {
		err := db.GetDB().Mysql.Table(conf.USER_TABLE_NAME).Create(userInfo).Error
		if err == nil {
			break
		}
		times++
	}
	if times == 10 {
		return errors.New("register failed")
	}

	userInfoBinary, err := json.Marshal(userInfo)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"err": err.Error(),
		}).Info("Register user info marshal failed.")
	}
	times = 0
	for times < 10 {
		setSuccess, err := u.DBInstance.Redis.HSet(conf.NAME_TO_USER_INFO_HASH_KEY, userName, userInfoBinary).Result()
		if setSuccess && err == nil {
			break
		}
	}
	if times == 10 {
		db.GetDB().Mysql.Table(conf.USER_TABLE_NAME).Delete(userInfo)
		return errors.New("register failed")
	}
	u.DBInstance.Redis.SAdd(conf.USER_NAME_SET_KEY, userName)

	return nil
}

func (u *UserDao) GetUserInfoWithUserName(ctx context.Context, userName string) (*model.UserItem, error) {
	userInfoBinary, err := u.DBInstance.Redis.HGet(conf.NAME_TO_USER_INFO_HASH_KEY, userName).Result()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"userName": userName,
			"err":      err.Error(),
		}).Info("GetUserInfoWithUserName failed.")
		return nil, err
	}

	userInfo := &model.UserItem{}
	err = json.Unmarshal([]byte(userInfoBinary), userInfo)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"userName": userName,
			"err":      err.Error(),
		}).Info("User info unmarshal failed.")
		return nil, err
	}

	return userInfo, nil
}

func (u *UserDao) GetUserInfoWithToken(ctx context.Context, token string) (*model.UserItem, error) {
	userInfoBinary, err := u.DBInstance.Redis.Get(token).Result()
	// Should always get userInfoBinary from redis
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"token": token,
			"err":   err.Error(),
		}).Info("GetUserInfoWithToken failed.")
		return nil, errors.New("GetUserInfoWithToken failed")
	}

	userInfo := &model.UserItem{}
	err = json.Unmarshal([]byte(userInfoBinary), userInfo)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"token": token,
			"err":   err.Error(),
		}).Info("User info unmarshal failed.")
		return nil, errors.New(err.Error())
	}

	return userInfo, err
}

func (u *UserDao) GetRandomEnemyUserInfo(ctx context.Context) (*model.UserItem, error) {
	randUserName, err := u.DBInstance.Redis.SRandMember(conf.USER_NAME_SET_KEY).Result()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"err": err.Error(),
		}).Info("Get rand user name failed.")
	}

	userInfo, err := u.GetUserInfoWithUserName(ctx, randUserName)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"randUserName": randUserName,
			"err":          err.Error(),
		}).Info("Get rand user info failed.")
	}

	return userInfo, nil
}

func (u *UserDao) availableWarehouseAddr(addr string) bool {
	if !strings.HasPrefix(addr, "git@github.com:") || !strings.HasSuffix(addr, "/GomokuGameImpl.git") {
		return false
	}
	return true
}

func (u *UserDao) setUserTokenPairInRedis(userInfo string) string {
	token := jwt.GetToken()
	u.DBInstance.Redis.Set(token, userInfo, conf.USER_INFO_TTL+1*time.Hour)
	return token
}
