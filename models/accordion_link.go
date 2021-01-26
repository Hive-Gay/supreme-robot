package models

import "time"

type Link struct {
	Title string `db:"title" ,json:"title"`
	Link  string `db:"link" ,json:"link"`

	ID        int       `db:"id" ,json:"id"`
	CreatedAt time.Time `db:"created_at" ,json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" ,json:"updated_at"`
}
