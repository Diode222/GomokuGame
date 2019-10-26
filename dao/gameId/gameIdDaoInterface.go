package gameId

import "context"

type GameIdDaoInterface interface {
	GetNextGameId(ctx context.Context) (int64, error)
}
