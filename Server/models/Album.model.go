package models


type Album struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Artist    string `json:"artist"`
	Year      int    `json:"year"`
	Cover_url string `json:"cover_url"`
}
