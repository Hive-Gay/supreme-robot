package models

import "time"

type Header struct {
	Title string `db:"title" ,json:"title"`

	LinkCount int `db:"link_count" ,json:"link_count"`

	ID        int `db:"id" ,json:"id"`
	CreatedAt time.Time `db:"created_at" ,json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" ,json:"updated_at"`
}

type Link struct {
	Title string `yaml:"title"`
	Link  string `yaml:"link"`
}


