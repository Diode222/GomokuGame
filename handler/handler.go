package handler

import (
	"GomokuGame/controller"
	"GomokuGame/dao/user"
	"GomokuGame/db"
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
			GameStartCtrl: controller.NewGameStartCtrl(),
		}
	})
	return h
}
