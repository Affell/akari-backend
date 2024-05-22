package user

import (
	"github.com/fatih/structs"
)

func (user User) ToSelfWebDetail() map[string]interface{} {
	m := structs.Map(user)
	delete(m, "password")
	delete(m, "reset_token")
	delete(m, "taxi_token")
	return m
}
