package authHandler

import (
	"akari/models"
	"akari/models/user"

	"github.com/kataras/golog"
	"github.com/kataras/iris/v12"
)

func Signout(c iris.Context, route models.Route) {

	if c.Method() != "POST" || route.Tertiary != "" || len(route.Tail) != 0 {
		c.StopWithStatus(iris.StatusNotFound)
		return
	}

	var signoutForm SignoutForm
	if err := c.ReadBody(&signoutForm); err != nil || len(signoutForm.Password) == 0 {
		c.StopWithJSON(iris.StatusBadRequest, iris.Map{"message": "Empty password"})
		return
	}

	var token user.UserToken
	if t := c.GetID(); t != nil {
		token = t.(user.UserToken)
	} else {
		c.StopWithJSON(iris.StatusUnauthorized, iris.Map{"message": "Invalid session"})
		return
	}

	if check := user.SecurityCheck(token.ID, signoutForm.Password); !check {
		c.StopWithJSON(iris.StatusForbidden, iris.Map{"message": "Invalid password"})
		return
	}

	if err := user.DeleteAccount(token.ID); err != "" {
		golog.Errorf("impossible de supprimer le compte %s : %s", token.ID, err)
		c.StopWithJSON(iris.StatusInternalServerError, iris.Map{"message": err})
		return
	}

	user.RevokeUserToken(token.TokenID)

	c.StopWithJSON(iris.StatusOK, iris.Map{"message": "You have successfully signed out of the website"})
}
