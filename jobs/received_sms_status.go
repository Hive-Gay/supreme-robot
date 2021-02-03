package jobs

import (
	"encoding/json"
	"github.com/gocraft/work"
)

const jobNameReceivedSMSStatus = "received_sms_status"

func (e *Enqueuer) ReceivedSMSStatus(params string) error {

	job, err := e.enqueuer.Enqueue(jobNameReceivedSMSStatus, work.Q{
		"params": params,
	})
	if err != nil {
		logger.Tracef("[%s] error enqueueing job: %s", jobNameReceivedSMSStatus, err.Error())
		return err
	}

	logger.Tracef("[%s](%s) job enqueued", jobNameReceivedSMSStatus, job.ID)
	return nil
}

func (c *Context) ReceivedSMSStatus(job *work.Job) error {
	// Extract arguments:
	params := job.ArgString("params")

	logger.Debugf("[%s](%s) processing sms status update", jobNameReceivedSMSStatus, job.ID)

	var smsData map[string][]string
	err := json.Unmarshal([]byte(params), &smsData)
	if err != nil {
		logger.Debugf("[%s](%s) can't unmarshal params: %s", jobNameReceivedSMSStatus, job.ID, err.Error())
		return err
	}

	sid := ""
	if val, ok := smsData["SmsSid"]; ok {
		if len(val) > 0 {
			sid = val[0]
		}
	}

	status := ""
	if val, ok := smsData["SmsStatus"]; ok {
		if len(val) > 0 {
			status = val[0]
		}
	}

	logger.Debugf("[%s](%s) sms status update: %s %s", jobNameReceivedSMSStatus, job.ID, sid, status)

	smsLog, err := c.modelclient.ReadSMSLogBySid(sid)
	if err != nil {
		logger.Errorf("[%s](%s) can't read sms log: %s", jobNameReceivedSMSStatus, job.ID, err.Error())
		return err
	}
	
	if smsLog == nil {
		logger.Warningf("[%s](%s) can't find sms log: %s", jobNameReceivedSMSStatus, job.ID, err.Error())
		return err
	}

	err = c.modelclient.UpdateSMSLogStatusBySid(smsLog, status)
	if err != nil {
		logger.Errorf("[%s](%s) can't update sms logs status: %s", jobNameReceivedSMSStatus, job.ID, err.Error())
		return err
	}

	return nil

}