package handler

import (
	"GomokuGame/controller"
	"GomokuGame/dao/gameId"
	"GomokuGame/dao/user"
	"GomokuGame/db"
	"GomokuGame/kube"
	"sync"
)

type handler struct {
	UserCtrl      *controller.UserCtrl
	GameStartCtrl *controller.GameStartCtrl
}

var h *handler
var handlerOnce sync.Once

func GetHandler() *handler {
	handlerOnce.Do(func() {
		h = &handler{
			UserCtrl: controller.NewUserCtrl(
				user.NewUserDao(
					db.GetDB(),
				),
			),
			GameStartCtrl: controller.NewGameStartCtrl(
				gameId.NewGameIdDao(db.GetDB().Redis),
				kube.NewKubePodsClient(),
			),
		}
	})
	return h
}
