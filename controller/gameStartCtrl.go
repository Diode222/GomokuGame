package controller

import (
	"GomokuGame/dao/gameId"
	"github.com/gin-gonic/gin"
	v1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

type GameStartCtrl struct {
	GameIdDao gameId.GameIdDaoInterface
	PodClient v1.PodInterface
}

func NewGameStartCtrl(gameIdDao gameId.GameIdDaoInterface, podClient v1.PodInterface) *GameStartCtrl {
	return &GameStartCtrl{
		GameIdDao: gameIdDao,
		PodClient: podClient,
	}
}

func (g *GameStartCtrl) Start(c *gin.Context) {

}
