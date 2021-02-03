package models

import (
	"database/sql"
	"time"
)

type SMSLog struct {
	AccountSid   string          `db:"account_sid"`
	ApiVersion   string          `db:"api_version"`
	Body         string          `db:"body"`
	Direction    string          `db:"direction"`
	ErrorCode    sql.NullInt32   `db:"error_code"`
	ErrorMessage sql.NullString  `db:"error_message"`
	FromID       int             `db:"from_id"`
	NumMedia     int             `db:"num_media"`
	NumSegments  int             `db:"num_segments"`
	Price        sql.NullFloat64 `db:"price"`
	PriceUnit    sql.NullString  `db:"price_unit"`
	Sid          string          `db:"sid"`
	Status       string          `db:"status"`
	ToID         int             `db:"to_id"`

	ID        int       `db:"id"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func (c *Client) CreateSMSLog(s *SMSLog) error {
	err := c.client.
		QueryRowx(`INSERT INTO public.sms_log(
			account_sid, api_version, body, direction, error_code, error_message, 
			from_id, num_media, num_segments, price, price_unit, sid, status, to_id)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
			RETURNING id, created_at, updated_at;`,
			s.AccountSid, s.ApiVersion, s.Body, s.Direction, s.ErrorCode, s.ErrorMessage,
			s.FromID, s.NumMedia, s.NumSegments, s.Price, s.PriceUnit, s.Sid, s.Status, s.ToID).
		Scan(&s.ID, &s.CreatedAt, &s.UpdatedAt)

	if err != nil {
		logger.Debugf("could not create record in sms_incoming_logL %s", err.Error())
	}

	return err
}

func (c *Client) ReadSMSLogBySid(sid string) (*SMSLog, error) {
	var sl SMSLog
	err := c.client.
		Get(&sl, `SELECT id, account_sid, api_version, body, direction, error_code, error_message, 
			from_id, num_media, num_segments, price, price_unit, sid, status, to_id, created_at, updated_at
			FROM public.sms_log WHERE sid = $1;`, sid)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &sl, nil
}

func (c *Client) UpdateSMSLogStatusBySid(sl *SMSLog, status string) error {
	err := c.client.
		QueryRowx(`UPDATE public.sms_log
			SET status=$2, updated_at=CURRENT_TIMESTAMP WHERE id=$1 RETURNING created_at, updated_at;`,
			sl.ID, status).
		Scan(&sl.CreatedAt, &sl.UpdatedAt)

	if err != nil {
		return err
	}

	sl.Status = status

	return nil
}