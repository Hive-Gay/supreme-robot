package models

import (
	"database/sql"
	"time"
)

type AccordionHeader struct {
	Title string `db:"title" ,json:"title"`

	LinkCount int `db:"link_count" ,json:"link_count"`

	ID        int       `db:"id" ,json:"id"`
	CreatedAt time.Time `db:"created_at" ,json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" ,json:"updated_at"`
}

func CountHiveHeaderLinks() (int, error) {
	var count int
	err := client.
		Get(&count, `SELECT count(id) FROM public.accordion_links WHERE accordion_header_id is NULL;`)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func CreateAccordionHeader(a *AccordionHeader) error {
	err := client.
		QueryRowx(`INSERT INTO public.accordion_headers(title) 
		VALUES ($1) RETURNING id, created_at, updated_at;`, a.Title).
		Scan(&a.ID, &a.CreatedAt, &a.UpdatedAt)

	return err
}

func DeleteAccordionHeader(id int) error {
	err := client.
		QueryRowx(`DELETE FROM accordion_headers WHERE id = $1;`, id).
		Scan()

	return err
}

func ReadAccordionHeader(id int) (*AccordionHeader, error) {
	var header AccordionHeader
	err := client.
		Get(&header, `SELECT id ,title, created_at, updated_at 
		FROM accordion_headers WHERE id = $1;`, id)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &header, nil
}

func ReadAccordionHeaders() ([]*AccordionHeader, error) {
	var ahs []*AccordionHeader
	err := client.
		Select(&ahs, `SELECT h.id ,h.title, h.created_at, h.updated_at, COUNT(l.id) as link_count
		FROM accordion_headers as h
		LEFT JOIN accordion_links as l ON h.id = l.accordion_header_id
		GROUP BY h.id
		ORDER BY title;`)
	if err != nil {
		return nil, err
	}

	return ahs, nil
}