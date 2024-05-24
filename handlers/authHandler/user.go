package authHandler

import (
	"akari/auth"
	"akari/models"
	"akari/models/user"

	"github.com/kataras/iris/v12"
)

func User(c iris.Context, route models.Route) {

	if route.Secondary != "" {
		c.StopWithStatus(iris.StatusNotFound)
		return
	}

	var token user.UserToken
	if t := c.GetID(); t != nil {
		token = t.(user.UserToken)
	} else {
		c.StopWithStatus(iris.StatusForbidden)
		return
	}

	var code int
	var data interface{}

	switch c.Method() {
	case "GET":
		u, err := user.GetUserById(token.ID)
		if err != nil {
			code, data = iris.StatusUnauthorized, iris.Map{"message": "Incorrect token"}
		} else {
			code, data = iris.StatusOK, iris.Map{"user": u.ToSelfWebDetail()}
		}
	case "POST":
		var editUserForm EditUserForm
		if err := c.ReadBody(&editUserForm); err != nil || (editUserForm.Email == "" && editUserForm.Username == "" && editUserForm.Password == "") {
			code, data = iris.StatusBadRequest, iris.Map{"message": "Please fill in at least one of the following fields : username, email, password"}
		} else {
			code, data = postUser(token, editUserForm)
		}
	default:
		code = iris.StatusNotFound
	}

	if data != nil {
		c.StopWithJSON(code, data)
	} else {
		c.StopWithStatus(code)
	}

}

func postUser(token user.UserToken, editUserForm EditUserForm) (code int, data interface{}) {
	if token.IsNil() {
		code = iris.StatusUnauthorized
		return
	}

	u, err := user.GetUserById(token.ID)
	if err != nil {
		code, data = iris.StatusNotFound, iris.Map{"message": err}
	}

	usernameChanged := editUserForm.Username != "" && editUserForm.Username != u.Username
	emailChanged := editUserForm.Email != "" && editUserForm.Email != u.Email
	passwordChanged := editUserForm.Password != "" && !user.SecurityCheck(token.ID, editUserForm.Password)

	if usernameChanged {
		if err := auth.ValidUsername(editUserForm.Username); err != "" {
			code, data = iris.StatusBadRequest, iris.Map{"message": err}
			return
		}

		if !user.CheckUsernameAvailability(editUserForm.Username) {
			code, data = iris.StatusConflict, iris.Map{"message": "Username not available"}
			return
		}

		u.Username = editUserForm.Username
	}

	if emailChanged {
		if !auth.ValidEmail(editUserForm.Email) {
			code, data = iris.StatusBadRequest, iris.Map{"message": "Invalid Email"}
			return
		}

		if !user.CheckEmailAvailability(editUserForm.Email) {
			code, data = iris.StatusConflict, iris.Map{"message": "Email not available"}
			return
		}
		u.Email = editUserForm.Email
	}

	if passwordChanged {
		if err := auth.ValidPassword(editUserForm.Password); err != "" {
			code, data = iris.StatusBadRequest, iris.Map{"message": err}
			return
		}

		u.Password = editUserForm.Password
	}

	if user.UpdateUser(u, passwordChanged) {
		token.Update(u)
		token.Store()
		code = iris.StatusOK
	} else {
		code = iris.StatusInternalServerError
	}

	return

}
