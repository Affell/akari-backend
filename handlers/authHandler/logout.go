package authHandler

import (
	"akari/models"
	"akari/models/user"

	"github.com/kataras/iris/v12"
)

func Logout(c iris.Context, route models.Route) {

	if c.Method() != "POST" || route.Tertiary != "" || len(route.Tail) != 0 {
		c.StopWithStatus(iris.StatusNotFound)
		return
	}

	var id string
	if t := c.GetID(); t != nil {
		id = t.(user.UserToken).TokenID
	} else {
		c.StopWithJSON(iris.StatusUnauthorized, iris.Map{"message": "Invalid session"})
		return
	}

	user.RevokeUserToken(id)
	c.StopWithStatus(iris.StatusOK)

}
