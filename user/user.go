package user

import (
	"GomokuGame/model"
	"github.com/gin-gonic/gin"
)

func Register(c *gin.Context) {
	userName := c.PostForm("user_name")
	password := c.PostForm("password")
	warehouseAddr := c.PostForm("warehouse_addr")

	user := &model.UserItem{
		UserName:      userName,
		Password:      password,
		WarehouseAddr: warehouseAddr,
	}

	c.String(user.Register())
}

func Login(c *gin.Context) {
	userName := c.PostForm("user_name")
	password := c.PostForm("password")

	user := &model.UserItem{
		UserName:      userName,
		Password:      password,
	}

	c.String(user.Login())
}
