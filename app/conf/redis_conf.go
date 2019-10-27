package conf

import "time"

const (
	REDIS_ADDR     string = "139.155.46.62:6379"
	REDIS_PASSWORD string = ""
	REDIS_DB_NAME  int    = 0

	USER_INFO_TTL time.Duration = 240 * time.Hour

	NAME_TO_USER_INFO_HASH_KEY = "hash_name_to_user_infos"

	USER_NAME_SET_KEY = "user_name_set_keys"

	GAME_RESULT_REDIS_STORE_TIME = 5 * time.Minute
)
