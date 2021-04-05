package cfip

import (
	"time"

	"github.com/gomodule/redigo/redis"
)

type LimiterConf struct {
	Namespace string

	Tokens   uint64
	Interval time.Duration

	Dialer func() (redis.Conn, error)
}
