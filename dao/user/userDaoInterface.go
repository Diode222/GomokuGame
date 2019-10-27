package user

import (
	"GomokuGame/model"
	"context"
)

type UserDaoInterface interface {
	Login(ctx context.Context, userName string, password string) (string, error)
	Register(ctx context.Context, userName string, password string, warehouseAddr string) error
	GetUserInfoWithUserName(ctx context.Context, userName string) (*model.UserItem, error)
	GetUserInfoWithToken(ctx context.Context, token string) (*model.UserItem, error)
	GetRandomEnemyUserInfo(ctx context.Context) (*model.UserItem, error)
}
