package path

import (
	"GomokuGame/app/conf"
	"path"
)

func GetGameLogDirPath(gameId string) string {
	return path.Join(conf.LOG_ROOT_PATH, gameId)
}
