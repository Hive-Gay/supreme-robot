package models

import "time"

type SMSIncomingLog struct {
	AccountSid    string `db:"account_sid" ,json:"account_sid"`
	ApiVersion    string `db:"api_version" ,json:"api_version"`
	Body          string `db:"body" ,json:"body"`
	Direction     string `db:"direction" ,json:"direction"`
	FromID        int    `db:"from_id" ,json:"from_id"`
	MessageSid    string `db:"message_sid" ,json:"message_sid"`
	NumMedia      int    `db:"num_media" ,json:"num_media"`
	NumSegments   int    `db:"num_segments" ,json:"num_segments"`
	SmsMessageSid string `db:"sms_message_sid" ,json:"sms_message_sid"`
	SmsSid        string `db:"sms_sid" ,json:"sms_sid"`
	SmsStatus     string `db:"sms_status" ,json:"sms_status"`
	ToID          int    `db:"to_id" ,json:"to_id"`

	ID        int       `db:"id" ,json:"id"`
	CreatedAt time.Time `db:"created_at" ,json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" ,json:"updated_at"`
}

func (c *Client) CreateSMSIncomingLog(s *SMSIncomingLog) error {
	err := c.client.
		QueryRowx(`INSERT INTO public.sms_incoming_log(
        account_sid, api_version, body, direction, from_id, message_sid, num_media, num_segments, sms_message_sid,
        sms_sid, sms_status, to_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		RETURNING id, created_at, updated_at;`,
			s.AccountSid, s.ApiVersion, s.Body, s.Direction, s.FromID, s.MessageSid, s.NumMedia, s.NumSegments,
			s.SmsMessageSid, s.SmsSid, s.SmsStatus, s.ToID).Scan(&s.ID, &s.CreatedAt, &s.UpdatedAt)

	if err != nil {
		logger.Debugf("could not create record in sms_incoming_logL %s", err.Error())
	}

	return err
}
