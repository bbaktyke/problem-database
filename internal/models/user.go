package models

type User struct {
	ID       int    `json:"-" db:"id"`
	Name     string `json:"name" validate:"required"`
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required,min=8"`
}

type Token struct {
	TokenString string `json:"token"`
}
