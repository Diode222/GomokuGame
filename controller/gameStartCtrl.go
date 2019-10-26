package controller

import "github.com/gin-gonic/gin"

type GameStartCtrl struct {
}

func NewGameStartCtrl() *GameStartCtrl {
	return &GameStartCtrl{}
}

func (g *GameStartCtrl) Start(c *gin.Context) {

}
