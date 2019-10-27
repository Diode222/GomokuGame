package router

import (
	"GomokuGame/handler"
	"GomokuGame/middleware"
	"GomokuGame/utils/api"
	"github.com/gin-gonic/gin"
	"net/http"
)

func InitRouter() {
	h := handler.GetHandler()

	router := gin.Default()
	router.NoRoute(api.NotFound)
	router.POST("/register", h.UserCtrl.Register)
	router.POST("/login", h.UserCtrl.Login)

	authorized := router.Group("/user", middleware.UserMiddleware())
	authorized.POST("/info", func(c *gin.Context) {
		c.String(http.StatusOK, "info")
	})

	gameRouter := router.Group("/game", middleware.UserMiddleware())
	gameRouter.POST("/start", h.GameStartCtrl.Start)

	router.Run(":8080")
}
