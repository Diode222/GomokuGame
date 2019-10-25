package handler

import (
	"GomokuGame/controller"
	"GomokuGame/dao/user"
	"GomokuGame/db"
	"sync"
)

type handler struct {
	UserCtrl *controller.UserCtrl
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
		}
	})
	return h
}
