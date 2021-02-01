package jobs

import "github.com/gocraft/work"

const jobNameReceivedSMS = "received_sms"
const jobKeyReceivedSMSDBID = "db_id"

func (e *Enqueuer) ReceivedSMS(dbID int) error {
	job, err := e.enqueuer.Enqueue(jobNameReceivedSMS, work.Q{
		jobKeyReceivedSMSDBID: dbID,
	})
	if err != nil {
		logger.Tracef("[%s] error enqueueing job: %s", jobNameReceivedSMS, err.Error())
		return err
	}

	logger.Tracef("[%s]($s) job enqueued", jobNameReceivedSMS, job.ID)
	return nil
}

func (c *Context) ReceivedSMS(job *work.Job) error {
	// Extract arguments:
	id := job.ArgString(jobKeyReceivedSMSDBID)

	logger.Debugf("[%s]($s) processing sms: %d", jobNameReceivedSMS, job.ID, id)

	return nil
}