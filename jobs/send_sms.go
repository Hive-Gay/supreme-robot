package jobs

import (
	"database/sql"
	"github.com/Hive-Gay/supreme-robot/models"
	"github.com/gocraft/work"
	"strconv"
	"time"
)

const jobNameSendSMS = "send_sms"

func (e *Enqueuer) SendSMS(fromID, toID int, body string) error {
	job, err := e.enqueuer.Enqueue(jobNameSendSMS, work.Q{
		"from_id": fromID,
		"to_id": toID,
		"body": body,
	})
	if err != nil {
		logger.Tracef("[%s] error enqueueing job: %s", jobNameSendSMS, err.Error())
		return err
	}

	logger.Tracef("[%s](%s) job enqueued", jobNameSendSMS, job.ID)
	return nil

	return nil
}

func (c *Context) SendSMS(job *work.Job) error {
	// Extract arguments:
	fromID := job.ArgInt64("from_id")
	toID := job.ArgInt64("to_id")
	body := job.ArgString("body")

	logger.Debugf("[%s](%s) sending sms", jobNameReceivedSMS, job.ID)

	from, err := c.modelclient.ReadPhoneNumber(int(fromID))
	if err != nil {
		logger.Debugf("[%s](%s) couldn't read phone number: %s", jobNameReceivedSMS, job.ID, err.Error())
		return err
	}

	to, err := c.modelclient.ReadPhoneNumber(int(toID))
	if err != nil {
		logger.Debugf("[%s](%s) couldn't read phone number: %s", jobNameReceivedSMS, job.ID, err.Error())
		return err
	}

	resp, err := c.twilioClient.SendMessage(from.Number, to.Number, body, "")
	if err != nil {
		logger.Debugf("[%s](%s) couldn't send message: %s", jobNameReceivedSMS, job.ID, err.Error())
		return err
	}

	numMedia, _ := strconv.Atoi(resp.NumMedia)
	numSegments, _ := strconv.Atoi(resp.NumSegments)

	dateCreated, err := time.Parse(time.RFC1123Z, resp.DateCreated)
	if err != nil {
		logger.Warningf("[%s](%s) can't parse time: %s", jobNameReceivedSMS, job.ID, err.Error())
	}

	dateUpdated, err := time.Parse(time.RFC1123Z, resp.DateUpdated)
	if err != nil {
		logger.Warningf("[%s](%s) can't parse time: %s", jobNameReceivedSMS, job.ID, err.Error())
	}

	smsLog := models.SMSLog{
		AccountSid: resp.AccountSid,
		ApiVersion: resp.ApiVersion,
		Body: resp.Body,
		DateCreated: sql.NullTime{Valid: true, Time: dateCreated},
		DateUpdated: sql.NullTime{Valid: true, Time: dateUpdated},
		Direction: resp.Direction,
		FromID:  from.ID,
		NumMedia: numMedia,
		NumSegments: numSegments,
		Sid: resp.Sid,
		Status: resp.Status,
		ToID: to.ID,
	}

	if resp.ErrorCode != nil {
		errorCode, _ := strconv.Atoi(*resp.ErrorCode)
		smsLog.ErrorCode = sql.NullInt32{
			Valid: true,
			Int32: int32(errorCode),
		}
	}

	if resp.ErrorMessage != nil {
		smsLog.ErrorMessage = sql.NullString{
			Valid: true,
			String: *resp.ErrorMessage,
		}
	}

	if resp.Price != nil {
		f, err := strconv.ParseFloat(*resp.Price, 64)
		if err != nil {
			logger.Warningf("[%s](%s) can't parse price: %s", jobNameReceivedSMS, job.ID, err.Error())
		}
		smsLog.Price = sql.NullFloat64{
			Valid: true,
			Float64: f,
		}
	}

	if resp.PriceUnit != nil {
		smsLog.PriceUnit = sql.NullString{
			Valid: true,
			String: *resp.PriceUnit,
		}
	}

	err = c.modelclient.CreateSMSIncomingLog(&smsLog)
	if err != nil {
		logger.Errorf("could not save sms: %s", err.Error())
		return err
	}

	return nil

}