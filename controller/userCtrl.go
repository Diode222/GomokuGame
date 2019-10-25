package controller

import (
	"GomokuGame/dao/user"
	"GomokuGame/utils/json"
	"github.com/gin-gonic/gin"
	"net/http"
)

type UserCtrl struct {
	UserDao user.UserDaoInterface
}

func NewUserCtrl(dao user.UserDaoInterface) *UserCtrl {
	return &UserCtrl{
		UserDao: dao,
	}
}

func (ctrl *UserCtrl) Login(c *gin.Context) {
	userName := c.PostForm("user_name")
	password := c.PostForm("password")

	token, err := ctrl.UserDao.Login(c.Request.Context(), userName, password)
	if err != nil {
		c.String(http.StatusUnauthorized, json.JsonResponse(http.StatusUnauthorized, "Login failed"))
	} else {
		c.String(http.StatusOK, json.JsonResponse(http.StatusOK, token))
	}
}

func (ctrl *UserCtrl) Register(c *gin.Context) {
	userName := c.PostForm("user_name")
	password := c.PostForm("password")
	warehouseAddr := c.PostForm("warehouse_addr")

	err := ctrl.UserDao.Register(c.Request.Context(), userName, password, warehouseAddr)
	if err != nil {
		c.String(http.StatusUnauthorized, json.JsonResponse(http.StatusUnauthorized, "register failed."))
	} else {
		c.String(http.StatusOK, json.JsonResponse(http.StatusOK, "register success."))
	}
}
