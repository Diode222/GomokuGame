package format

import (
	"GomokuGame/model"
	"encoding/json"
	"github.com/sirupsen/logrus"
)

func GameResultFormatter(gameResult *model.MatchResultItem) *GameResultFormat {
	if gameResult == nil {
		logrus.Info("Game result is nil.")
		return &GameResultFormat{
			ServerError: true,
		}
	}

	// This match is still going
	if gameResult.GameStatus {
		return &GameResultFormat{
			GameID:      gameResult.GameID,
			GameStatus:  true,
			ServerError: false,
		}
	}

	operationsJson := gameResult.Operations
	operationsModel := []*model.Operation{}
	if err := json.Unmarshal([]byte(operationsJson), &operationsModel); err != nil {
		logrus.WithFields(logrus.Fields{
			"GameId": gameResult.GameID,
			"err":    err.Error(),
		}).Info("Unmarshal operations of game result failed.")
		return &GameResultFormat{
			GameID:      gameResult.GameID,
			StartTime:   gameResult.StartTime,
			EndTime:     gameResult.EndTime,
			ServerError: true,
		}
	}

	gameOperationsFormat := []*GameOperationFormat{}
	for _, operationModel := range operationsModel {
		if operationsModel == nil || operationModel.Position == nil {
			continue
		}

		gameOperationsFormat = append(gameOperationsFormat, &GameOperationFormat{
			Player:    operationModel.Player,
			PositionX: int(operationModel.Position.X),
			PositionY: int(operationModel.Position.Y),
			Type:      operationModel.Type,
		})
	}

	return &GameResultFormat{
		GameID:           gameResult.GameID,
		GameStatus:       false,
		BoardLength:      gameResult.BoardLength,
		BoardHeight:      gameResult.BoardHeight,
		Player1ID:        gameResult.Player1ID,
		Player2ID:        gameResult.Player2ID,
		Player1FirstHand: gameResult.Player1FirstHand,
		MaxThinkingTime:  gameResult.MaxThinkingTime,
		Winner:           gameResult.Winner,
		StartTime:        gameResult.StartTime,
		EndTime:          gameResult.EndTime,
		Operations:       gameOperationsFormat,
		FoulPlayer:       gameResult.FoulPlayer,
		ServerError:      gameResult.ServerError,
	}
}
