package gameResult

import (
	"GomokuGame/app/conf"
	"GomokuGame/db"
	"GomokuGame/model"
	"context"
	"encoding/json"
	"errors"
	"github.com/sirupsen/logrus"
)

type GameResultDao struct {
	dbInstance *db.DB
}

const (
	GAME_IS_GOING_FLAG     = "gaming"
	GAME_RESULT_KEY_PREFIX = "game_result_"
)

func NewGameResultDao(instance *db.DB) *GameResultDao {
	return &GameResultDao{
		dbInstance: instance,
	}
}

func (g *GameResultDao) SetTempGamingStatusInRedis(ctx context.Context, gameId string) {
	g.dbInstance.Redis.Set(GAME_RESULT_KEY_PREFIX+gameId, GAME_IS_GOING_FLAG, conf.MAX_PULL_DATA_TIME)
}

func (g *GameResultDao) GetGameResult(ctx context.Context, gameId string) (*model.MatchResultItem, error) {
	gameResultJson, err := g.dbInstance.Redis.Get(GAME_RESULT_KEY_PREFIX + gameId).Result()

	if gameResultJson == GAME_IS_GOING_FLAG && err == nil {
		return &model.MatchResultItem{
			GameID:      gameId,
			GameStatus:  true,
			ServerError: false,
		}, nil
	}

	gameResult := &model.MatchResultItem{}
	if gameResultJson != "" && err == nil {
		unmarshalErr := json.Unmarshal([]byte(gameResultJson), gameResult)
		if unmarshalErr == nil {
			return gameResult, nil
		}
	}

	g.dbInstance.Mysql.Table(conf.MATCH_RESULT_TABLE_NAME).Where("game_id = ?", gameId).Find(gameResult)
	if gameResult.GameID == "" {
		logrus.WithFields(logrus.Fields{
			"gameId": gameId,
		}).Info("GetGameResult can not find match.")
		return nil, errors.New("GetGameResult can not find match.")
	}

	j, err := json.Marshal(gameResult)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"gameId": gameId,
			"err":    err.Error(),
		}).Info("GetGameResult marshal game result model failed.")
		return nil, errors.New("GetGameResult marshal game result model failed.")
	}
	g.dbInstance.Redis.Set(GAME_RESULT_KEY_PREFIX+gameId, j, conf.GAME_RESULT_REDIS_STORE_TIME)

	return gameResult, nil
}
