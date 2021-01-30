package models

import (
	"database/sql"
	"time"
)

type AccordionLink struct {
	Title string `db:"title" ,json:"title"`
	Link  string `db:"link" ,json:"link"`

	AccordionHeaderID sql.NullInt32 `db:"accordion_header_id"`

	ID        int       `db:"id" ,json:"id"`
	CreatedAt time.Time `db:"created_at" ,json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" ,json:"updated_at"`
}

func CreateAccordionLink(a *AccordionLink) error {
	err := client.
		QueryRowx(`INSERT INTO public.accordion_links(accordion_header_id, title, link) 
		VALUES ($1, $2, $3) RETURNING id, created_at, updated_at;`, a.AccordionHeaderID, a.Title, a.Link).
		Scan(&a.ID, &a.CreatedAt, &a.UpdatedAt)

	return err
}

func DeleteAccordionLink(id int) error {
	err := client.
		QueryRowx(`DELETE FROM accordion_links WHERE id = $1;`, id).
		Scan()
	if err == sql.ErrNoRows {
		return nil
	}

	return err
}

func ReadAccordionLinks(headerID sql.NullInt32) ([]*AccordionLink, error) {
	var als []*AccordionLink

	var err error
	if headerID.Valid {
		err = client.
			Select(&als, `SELECT id, accordion_header_id, title, link, created_at, updated_at 
			FROM accordion_links WHERE accordion_header_id = $1 ORDER BY title;`, headerID.Int32)
	} else {
		err = client.
			Select(&als, `SELECT id, accordion_header_id, title, link, created_at, updated_at 
			FROM accordion_links WHERE accordion_header_id is NULL ORDER BY title;`)
	}
	if err != nil {
		return nil, err
	}

	return als, nil
}


func ReadAccordionLink(headerID sql.NullInt32, linkID int) (*AccordionLink, error) {
	var link AccordionLink

	var err error
	if headerID.Valid {
		logger.Tracef("getting header %d link %d", headerID.Int32, linkID)
		err = client.
			Get(&link, `SELECT id, accordion_header_id, title, link, created_at, updated_at 
			FROM accordion_links WHERE id = $1 AND accordion_header_id = $2;`, linkID, headerID.Int32)
	} else {
		logger.Tracef("getting header NULL link %d", linkID)
		err = client.
			Get(&link, `SELECT id, accordion_header_id, title, link, created_at, updated_at 
			FROM accordion_links WHERE id = $1 AND accordion_header_id is NULL;`, linkID)
	}

	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &link, nil
}

func UpdateAccordionLink(a *AccordionLink) error {
	err := client.
		QueryRowx(`UPDATE public.accordion_links
		SET title=$1, link=$2, updated_at=CURRENT_TIMESTAMP
		WHERE id=$3 RETURNING updated_at;`, a.Title, a.Link, a.ID).
		Scan(&a.UpdatedAt)

	return err
}