package jobs

import (
	"github.com/garyburd/redigo/redis"
	"github.com/juju/loggo"
)

var logger = loggo.GetLogger("jobs")

type Enqueuer struct {
	namespace string
	pool      *redis.Pool
}

func NewEnqueuer(namespace string, pool *redis.Pool) *Enqueuer {
	logger.Debugf("creating new enqueuer in namespace %s", namespace)
	return &Enqueuer{
		namespace: namespace,
		pool:      pool,
	}
}

func (e *Enqueuer)ReceivedSMS() {

}