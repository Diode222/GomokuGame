package router

import (
	"GomokuGame/middleware"
	"GomokuGame/user"
	"GomokuGame/utils/api"
	"github.com/gin-gonic/gin"
	"net/http"
)

func InitRouter() {
	router := gin.Default()
	router.NoRoute(api.NotFound)
	router.POST("/register", user.Register)
	router.POST("/login", user.Login)

	authorized := router.Group("/user", middleware.UserMiddleware())
	authorized.POST("/info", func(c *gin.Context) {
		c.String(http.StatusOK, "info")
	})

	router.Run(":8080")
}
