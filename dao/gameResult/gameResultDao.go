package gameResult

import (
	"GomokuGame/app/conf"
	"GomokuGame/db"
	"context"
)

type GameResultDao struct {
	dbInstance *db.DB
}

func NewGameResultDao(instance *db.DB) *GameResultDao {
	return &GameResultDao{
		dbInstance: instance,
	}
}

func (g *GameResultDao) SetTempGamingStatusInRedis(ctx context.Context, gameId string) {
	g.dbInstance.Redis.Set("game_result_"+gameId, "gaming", conf.MAX_PULL_DATA_TIME)
}
