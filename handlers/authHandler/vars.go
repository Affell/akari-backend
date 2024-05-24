package authHandler

import (
	"akari/models"
)

const (
	Service string = "auth"
)

var (
	Handlers models.HandlerMap = models.HandlerMap{
		"signup":  Signup,
		"signout": Signout,
		"login":   Login,
		"logout":  Logout,
		"user":    User,
	}
)

type (
	SignupForm struct {
		Username string `form:"username" json:"username" binding:"required"`
		Email    string `form:"email" json:"email" binding:"required"`
		Password string `form:"password" json:"password" binding:"required"`
	}

	LoginForm struct {
		Username string `form:"username" json:"username"`
		Password string `form:"password" json:"password"`
		Token    string `form:"token" json:"token"`
	}

	SignoutForm struct {
		Password string `form:"password" json:"password" binding:"required"`
	}

	EditUserForm struct {
		Email    string `structs:"email" form:"email" json:"email" binding:"required"`
		Username string `structs:"username" form:"username" json:"username" binding:"required"`
		Password string `structs:"password,omitempty" json:"password"`
	}

	AskRecoverForm struct {
		Email string `form:"email" json:"email" binding:"required"`
	}

	RecoverForm struct {
		Password string `form:"password" json:"password" binding:"required"`
	}

	PermissionQueries struct {
		Permission string `form:"permission" json:"permission" binding:"required"`
	}

	ExchangeDiscordForm struct {
		Code string
	}

	VerifyEmailForm struct {
		Token string
	}
)
