package models

import "database/sql"

type Album struct {
	ID        int            `json:"id"`
	Title     string         `json:"title"`
	Artist    string         `json:"artist"`
	Year      int            `json:"year"`
	Cover_url sql.NullString `json:"cover_url"`
}
