package models

type User struct {
	ID       int    `json:"id"`
	FullName string `json:"fullname"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Input struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
