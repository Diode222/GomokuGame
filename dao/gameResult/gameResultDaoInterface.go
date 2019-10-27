package gameResult

import (
	"GomokuGame/model"
	"context"
)

type GameResultDaoInterface interface {
	SetTempGamingStatusInRedis(ctx context.Context, gameId string)
	GetGameResult(ctx context.Context, gameId string) (*model.MatchResultItem, error)
}
