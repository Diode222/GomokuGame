package handler

import (
	"GomokuGame/controller"
	"GomokuGame/dao/gameId"
	"GomokuGame/dao/gameResult"
	"GomokuGame/dao/user"
	"GomokuGame/db"
	"GomokuGame/kube"
	"sync"
)

type handler struct {
	UserCtrl *controller.UserCtrl
	GameCtrl *controller.GameCtrl
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
			GameCtrl: controller.NewGameCtrl(
				gameId.NewGameIdDao(db.GetDB().Redis),
				user.NewUserDao(db.GetDB()),
				gameResult.NewGameResultDao(db.GetDB()),
				kube.NewKubePodsClient(),
			),
		}
	})
	return h
}
