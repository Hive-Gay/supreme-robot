package jobs

import (
	"github.com/gocraft/work"
	"github.com/gomodule/redigo/redis"
	"github.com/juju/loggo"
)

var logger = loggo.GetLogger("jobs")

type Enqueuer struct {
	enqueuer *work.Enqueuer
}

func NewEnqueuer(namespace string, redisAddress string) *Enqueuer {
	logger.Debugf("creating new enqueuer in namespace %s", namespace)

	var enqueuer = work.NewEnqueuer(namespace, &redis.Pool{
		MaxActive: 5,
		MaxIdle: 5,
		Wait: true,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", redisAddress)
		},
	})

	return &Enqueuer{
		enqueuer: enqueuer,
	}
}

