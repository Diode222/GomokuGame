package conf

import "time"

const (
	NSQD_PORT       = 4151
	NSQLOOKUPD_PORT = 4161

	MAX_PULL_DATA_TIME = 5 * time.Minute
)

var (
	NSQD_ADDR  string = ""
	NSQLOOKUPD string = ""
)
