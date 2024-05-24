package auth

import (
	"akari/models/user"

	"github.com/kataras/iris/v12"
)

const (
	TokenKeyName = "INSAkari-Connect-Token"
)

type Header struct {
	TokenID string `header:"INSAkari-Connect-Token"`
}

func AuthRequired() func(iris.Context) {
	return func(c iris.Context) {

		var header Header
		if err := c.ReadHeaders(&header); err != nil || len(header.TokenID) != 36 {
			c.Next()
			return
		}

		Token, err := user.GetUserToken(header.TokenID)
		if err == nil {
			c.SetID(Token)
		}

		c.Next()
	}
}
