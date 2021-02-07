package jobs

import (
	"database/sql"
	"encoding/json"
	"github.com/Hive-Gay/supreme-robot/models"
	"github.com/gocraft/work"
	"strconv"
)

const jobNameReceivedSMS = "received_sms"

func (e *Enqueuer) ReceivedSMS(params string) error {

	job, err := e.enqueuer.Enqueue(jobNameReceivedSMS, work.Q{
		"params": params,
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
	params := job.ArgString("params")

	logger.Debugf("[%s](%s) processing sms", jobNameReceivedSMS, job.ID)
	logger.Tracef("[%s](%s) params: %s", params)

	var smsData map[string][]string
	err := json.Unmarshal([]byte(params), &smsData)
	if err != nil {
		logger.Debugf("[%s](%s) couldn't unmarshal params: %s", jobNameReceivedSMS, job.ID, err.Error())
		return err
	}

	logger.Tracef("[%s](%s) smsdata : %#v", jobNameReceivedSMS, job.ID, smsData["AccountSid"][0])

	// Get From Object
	fromPhoneNumber := models.PhoneNumber{
		City: sql.NullString{Valid: false},
		Country: sql.NullString{Valid: false},
		State: sql.NullString{Valid: false},
		Zip: sql.NullString{Valid: false},
	}

	if val, ok := smsData["From"]; ok {
		if len(val) > 0 {
			fromPhoneNumber.Number = val[0]
		}
	}

	if val, ok := smsData["FromCity"]; ok {
		if len(val) > 0 && val[0] != ""{
			fromPhoneNumber.City = sql.NullString{
				Valid: true,
				String: val[0],
			}
		}
	}

	if val, ok := smsData["FromCountry"]; ok {
		if len(val) > 0 && val[0] != ""{
			fromPhoneNumber.Country = sql.NullString{
				Valid: true,
				String: val[0],
			}
		}
	}

	if val, ok := smsData["FromState"]; ok {
		if len(val) > 0 && val[0] != ""{
			fromPhoneNumber.State = sql.NullString{
				Valid: true,
				String: val[0],
			}
		}
	}

	if val, ok := smsData["FromZip"]; ok {
		if len(val) > 0 && val[0] != ""{
			fromPhoneNumber.Zip = sql.NullString{
				Valid: true,
				String: val[0],
			}
		}
	}

	err = fromPhoneNumber.Upsert(c.modelclient)
	if err != nil {
		logger.Debugf("[%s](%s) couldn't read number from database: %s", jobNameReceivedSMS, job.ID, err.Error())
		return err
	}

	// Get To Object
	toPhoneNumber := models.PhoneNumber{
		City: sql.NullString{Valid: false},
		Country: sql.NullString{Valid: false},
		State: sql.NullString{Valid: false},
		Zip: sql.NullString{Valid: false},
	}

	if val, ok := smsData["To"]; ok {
		if len(val) > 0 {
			toPhoneNumber.Number = val[0]
		}
	}

	if val, ok := smsData["ToCity"]; ok {
		if len(val) > 0 && val[0] != ""{
			toPhoneNumber.City = sql.NullString{
				Valid: true,
				String: val[0],
			}
		}
	}

	if val, ok := smsData["ToCountry"]; ok {
		if len(val) > 0 && val[0] != ""{
			toPhoneNumber.Country = sql.NullString{
				Valid: true,
				String: val[0],
			}
		}
	}

	if val, ok := smsData["ToState"]; ok {
		if len(val) > 0 && val[0] != ""{
			toPhoneNumber.State = sql.NullString{
				Valid: true,
				String: val[0],
			}
		}
	}

	if val, ok := smsData["ToZip"]; ok {
		if len(val) > 0 && val[0] != ""{
			toPhoneNumber.Zip = sql.NullString{
				Valid: true,
				String: val[0],
			}
		}
	}

	err = toPhoneNumber.Upsert(c.modelclient)
	if err != nil {
		logger.Debugf("[%s](%s) couldn't read number To database: %s", jobNameReceivedSMS, job.ID, err.Error())
		return err
	}

	// Populate SMS
	smsLog := models.SMSConversationLine{
		Direction: "incoming",
		FromID: fromPhoneNumber.ID,
		ToID: toPhoneNumber.ID,
	}

	if val, ok := smsData["AccountSid"]; ok {
		if len(val) > 0 {
			smsLog.AccountSid = val[0]
		}
	}

	if val, ok := smsData["ApiVersion"]; ok {
		if len(val) > 0 {
			smsLog.ApiVersion = val[0]
		}
	}

	if val, ok := smsData["Body"]; ok {
		if len(val) > 0 {
			smsLog.Body = val[0]
		}
	}

	if val, ok := smsData["NumMedia"]; ok {
		if len(val) > 0 {
			i, err := strconv.Atoi(val[0])
			if err == nil {
				smsLog.NumMedia = i
			}
		}
	}

	if val, ok := smsData["NumSegments"]; ok {
		if len(val) > 0 {
			i, err := strconv.Atoi(val[0])
			if err == nil {
				smsLog.NumSegments = i
			}
		}
	}

	if val, ok := smsData["SmsSid"]; ok {
		if len(val) > 0 {
			smsLog.Sid = val[0]
		}
	}

	if val, ok := smsData["SmsStatus"]; ok {
		if len(val) > 0 {
			smsLog.Status = val[0]
		}
	}

	logger.Tracef("[%s](%s) writing sms to database", jobNameReceivedSMS, job.ID)
	err = smsLog.Create(c.modelclient)

	return err
}