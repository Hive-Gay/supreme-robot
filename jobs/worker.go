package jobs

import (
	"github.com/Hive-Gay/supreme-robot/models"
	"github.com/gocraft/work"
	"github.com/gomodule/redigo/redis"
)

type Worker struct {
	modelclient *models.Client
	pool        *work.WorkerPool
}

func NewWorker(namespace string, redisAddress string, mc *models.Client) *Worker {
	logger.Debugf("creating new worker in namespace %s", namespace)

	var pool = work.NewWorkerPool(Context{}, 10, namespace, &redis.Pool{
		MaxActive: 5,
		MaxIdle:   5,
		Wait:      true,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", redisAddress)
		},
	})

	// Map the name of jobs to handler functions
	pool.Job(jobNameReceivedSMS, (*Context).ReceivedSMS)

	return &Worker{
		modelclient: mc,
		pool: pool,
	}
}

func (w *Worker) Start() {
	w.pool.Start()
}

func (w *Worker) Stop() {
	w.pool.Stop()
}
