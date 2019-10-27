package conf

import "time"

const (
	NSQD_PORT       = ":4151"
	NSQLOOKUPD_PORT = ":4161"
	NSQ_PUB_PORT    = ":4150"

	MAX_PULL_DATA_TIME = 5 * time.Minute
)

var (
	NSQ_PUB_ADDR   string = "10.244.0.1" + NSQ_PUB_PORT
	NSQD_PULL_ADDR string = "10.244.0.1" + NSQD_PORT
	NSQLOOKUPD_ADDR     string = "10.244.0.2" + NSQLOOKUPD_PORT
)
