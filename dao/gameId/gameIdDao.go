package gameId

import (
	"context"
	"github.com/go-redis/redis"
)

type GameIdDao struct {
	redisInstance *redis.Client
}

func NewGameIdDao(instance *redis.Client) *GameIdDao {
	return &GameIdDao{
		redisInstance: instance,
	}
}

func (g *GameIdDao) GetNextGameId(ctx context.Context) (int64, error) {
	return g.redisInstance.Incr("global_game_id").Result()
}
