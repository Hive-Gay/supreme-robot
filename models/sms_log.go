package models

import (
	"database/sql"
	"time"
)

type SMSLog struct {
	AccountSid   string          `db:"account_sid"`
	ApiVersion   string          `db:"api_version"`
	Body         string          `db:"body"`
	DateCreated  sql.NullTime    `db:"date_created"`
	DateSent     sql.NullTime    `db:"date_sent"`
	DateUpdated  sql.NullTime    `db:"date_updated"`
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

func (c *Client) CreateSMSIncomingLog(s *SMSLog) error {
	err := c.client.
		QueryRowx(`INSERT INTO public.sms_log(
        account_sid, api_version, body, date_created, date_sent, date_updated, direction, error_code, error_message, 
        from_id, num_media, num_segments, price, price_unit, sid, status, to_id)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17)
		RETURNING id, created_at, updated_at;`,
			s.AccountSid, s.ApiVersion, s.Body, s.DateCreated, s.DateSent, s.DateUpdated, s.Direction, s.ErrorCode,
			s.ErrorMessage, s.FromID, s.NumMedia, s.NumSegments, s.Price, s.PriceUnit, s.Sid, s.Status, s.ToID).
			Scan(&s.ID, &s.CreatedAt, &s.UpdatedAt)

	if err != nil {
		logger.Debugf("could not create record in sms_incoming_logL %s", err.Error())
	}

	return err
}
