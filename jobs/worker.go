package jobs

import (
	"github.com/Hive-Gay/supreme-robot/models"
	"github.com/Hive-Gay/supreme-robot/twilio"
	"github.com/gocraft/work"
	"github.com/gomodule/redigo/redis"
)

type Context struct{
	modelclient *models.Client
	twilioClient *twilio.Client
	webappHostname string
}

type Worker struct {
	pool        *work.WorkerPool
}

func NewWorker(namespace, webappHostname, redisAddress string, mc *models.Client, tc *twilio.Client) *Worker {
	logger.Debugf("creating new worker in namespace %s", namespace)

	var pool = work.NewWorkerPool(Context{}, 10, namespace, &redis.Pool{
		MaxActive: 5,
		MaxIdle:   5,
		Wait:      true,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", redisAddress)
		},
	})

	c := Context{
		modelclient: mc,
		twilioClient: tc,
		webappHostname: webappHostname,
	}

	// Map the name of jobs to handler functions
	pool.Job(jobNameReceivedSMS, c.ReceivedSMS)
	pool.Job(jobNameReceivedSMSStatus, c.ReceivedSMSStatus)
	pool.Job(jobNameSendSMS, c.SendSMS)

	return &Worker{
		pool: pool,
	}
}

func (w *Worker) Start() {
	w.pool.Start()
}

func (w *Worker) Stop() {
	w.pool.Stop()
}
