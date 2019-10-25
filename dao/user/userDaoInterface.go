package user

import "context"

type UserDaoInterface interface {
	Login(ctx context.Context, userName string, password string) (string, error)
	Register(ctx context.Context, userName string, password string, warehouseAddr string) error
}
