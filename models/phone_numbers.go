package models

import (
	"database/sql"
	"time"
)

type PhoneNumber struct {
	Number  string `db:"num" ,json:"num"`
	City    sql.NullString `db:"city" ,json:"city"`
	Country sql.NullString `db:"country" ,json:"country"`
	State   sql.NullString `db:"state" ,json:"state"`
	Zip     sql.NullString `db:"zip" ,json:"zip"`

	ID        int       `db:"id" ,json:"id"`
	CreatedAt time.Time `db:"created_at" ,json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" ,json:"updated_at"`
}

func (c *Client)CreatePhoneNumber(a *PhoneNumber) error {
	err := c.client.
		QueryRowx(`INSERT INTO public.phone_numbers(num, city, country, state, zip) 
		VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at, updated_at;`, a.Number, a.City, a.Country, a.State, a.Zip).
		Scan(&a.ID, &a.CreatedAt, &a.UpdatedAt)

	return err
}

func (c *Client)ReadPhoneNumber(id int) (*PhoneNumber, error) {
	var pn PhoneNumber
	err := c.client.
		Get(&pn, `SELECT id, num, city, country, state, zip, created_at, updated_at
		FROM public.phone_numbers WHERE id = $1;`, id)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &pn, nil
}

func (c *Client)ReadPhoneNumberByNumber(num string) (*PhoneNumber, error) {
	var pn PhoneNumber
	err := c.client.
		Get(&pn, `SELECT id, num, city, country, state, zip, created_at, updated_at
		FROM public.phone_numbers WHERE num = $1;`, num)
		if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &pn, nil
}

func (c *Client)UpdatePhoneNumber(pn *PhoneNumber) error {
	err := c.client.
		QueryRowx(`UPDATE public.phone_numbers
		SET city=$2, country=$3, state=$4, zip=$5, updated_at=CURRENT_TIMESTAMP
		WHERE id=$1 RETURNING created_at, updated_at;`, pn.ID, pn.City, pn.Country, pn.State, pn.Zip).
		Scan(&pn.CreatedAt, &pn.UpdatedAt)

	return err
}

func (c *Client)UpsertPhoneNumber(pn *PhoneNumber) error {
	var foundPN *PhoneNumber
	var err error

	// Check for existing record
	foundPN, err = c.ReadPhoneNumberByNumber(pn.Number)
	if err != nil {
		logger.Debugf("upsert could not read phone number: %s", err.Error())
		return err
	}

	// if not exist, create
	if foundPN == nil {
		err = c.CreatePhoneNumber(pn)
		if err != nil {
			logger.Debugf("upsert could not create phone number: %s", err.Error())
			return err
		}
		return nil
	}

	// look for changes
	changes := false

	if pn.City.Valid != foundPN.City.Valid || pn.City.String != foundPN.City.String {
		changes = true
	}
	if pn.Country.Valid != foundPN.Country.Valid || pn.Country.String != foundPN.Country.String {
		changes = true
	}
	if pn.State.Valid != foundPN.State.Valid || pn.State.String != foundPN.State.String {
		changes = true
	}
	if pn.Zip.Valid != foundPN.Zip.Valid || pn.Zip.String != foundPN.Zip.String {
		changes = true
	}

	// if changes update
	if changes {
		logger.Debugf("upsert is updating the metadata for: %s", pn.Number)
		pn.ID = foundPN.ID
		err = c.UpdatePhoneNumber(pn)
		if err != nil {
			logger.Debugf("upsert could not update phone number: %s", err.Error())
			return err
		}
		return nil
	}

	pn.ID = foundPN.ID
	pn.CreatedAt = foundPN.CreatedAt
	pn.UpdatedAt = foundPN.UpdatedAt
	return nil
}