package jobs

import (
	"github.com/Hive-Gay/supreme-robot/models"
	"github.com/gocraft/work"
	"github.com/google/uuid"
)

const jobNameReceivedSMS = "received_sms"

func (e *Enqueuer) ReceivedSMS(url, params, iToken string) error {

	job, err := e.enqueuer.Enqueue(jobNameReceivedSMS, work.Q{
		"url": url,
		"params": params,
		"i_token": iToken,
	})
	if err != nil {
		logger.Tracef("[%s] error enqueueing job: %s", jobNameReceivedSMS, err.Error())
		return err
	}

	logger.Tracef("[%s](%s) job enqueued", jobNameReceivedSMS, job.ID)
	return nil
}

func (c *Context) ReceivedSMS(job *work.Job) error {
	// Extract arguments:
	url := job.ArgString("url")
	params := job.ArgString("params")
	iToken, err := uuid.Parse(job.ArgString("i_token"))
	if err != nil {
		logger.Debugf("couldn't parse uuid: %s", err.Error())
		return err
	}

	logger.Debugf("[%s](%s) processing sms: %s", jobNameReceivedSMS, job.ID, iToken)

	smsLog := &models.SMSWebhookLog{
		URL: url,
		Params: params,
		IToken: iToken,
	}

	logger.Tracef("smsLog: %#v", smsLog)


	return nil
}