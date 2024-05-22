package user

import (
	"time"
)

const PERMISSIONS_EDIT_PERMISSION = "edit.user.permission"

var SensibleFields []string = []string{
	"email",
	"password",
}

type UserToken struct {
	TokenID   string    `json:"token_id" structs:"-"`
	ID        int64     `json:"id" structs:"id"`
	Email     string    `json:"email" structs:"email"`
	CreatedAt time.Time `json:"created_at" structs:"-"`
}

type User struct {
	ID         int64  `structs:"id"`
	Firstname  string `structs:"firstname"`
	Lastname   string `structs:"lastname"`
	Email      string `structs:"email"`
	Password   string `structs:"password"`
	TaxiToken  string `structs:"taxi_token"`
	ResetToken string `structs:"reset_token"`
}
