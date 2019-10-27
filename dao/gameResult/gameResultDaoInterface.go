package gameResult

import "context"

type GameResultDaoInterface interface {
	SetTempGamingStatusInRedis(ctx context.Context, gameId string)
}
