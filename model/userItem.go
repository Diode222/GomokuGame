package model

import (
	"GomokuGame/conf"
	"GomokuGame/db"
	"GomokuGame/utils/json"
	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	"net/http"
	"strings"
	"time"
)

type UserItem struct {
	gorm.Model
	UserName      string	`gorm:"not null"`
	Password      string	`gorm:"not null"`
	WarehouseAddr string	`gorm:"not null"`
}

// Need redis here, read from redis, save into mysql
func (u *UserItem) Register() (int, string) {
	if !availableWarehouseAddr(u.WarehouseAddr) {
		return http.StatusBadRequest, json.JsonResponse(http.StatusBadRequest, "Git warehouse address is not right.")
	}

	users := []*UserItem{}
	db.GetDB().Mysql.Table(conf.USER_TABLE_NAME).Where("user_name = ?", u.UserName).Where("password = ?", u.Password).Find(&users)
	if len(users) > 0 {
		return http.StatusCreated, json.JsonResponse(http.StatusCreated, "Username has registered.")
	}

	db.GetDB().Mysql.Table(conf.USER_TABLE_NAME).Create(&UserItem{
		UserName:      u.UserName,
		Password:      u.Password,
		WarehouseAddr: u.WarehouseAddr,
	})
	return http.StatusOK, json.JsonResponse(http.StatusOK, "Register success.")
}

func (u *UserItem) Login() (int, string) {
	user := &UserItem{}
	db.GetDB().Mysql.Table(conf.USER_TABLE_NAME).Where("user_name = ?", u.UserName).Find(user)
	if user.UserName != "" && user.Password != "" {
		return http.StatusOK, json.JsonResponse(http.StatusOK, setToken())
	} else {
		return http.StatusUnauthorized, json.JsonResponse(http.StatusUnauthorized, "Login failed.")
	}
}

func setToken() string {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := make(jwt.MapClaims)
	claims["exp"] = time.Now().Add(time.Hour * time.Duration(240)).Unix()
	claims["iat"] = time.Now().Unix()
	token.Claims = claims

	tokenString, err := token.SignedString([]byte("webapp"))
	if err != nil {
		return ""
	}
	return tokenString
}

func availableWarehouseAddr(addr string) bool {
	if !strings.HasPrefix(addr, "git@github.com:") || !strings.HasSuffix(addr, "/GomokuGameImpl.git") {
		return false
	}
	return true
}
