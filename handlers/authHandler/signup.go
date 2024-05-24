package authHandler

import (
	"akari/auth"
	"akari/models"
	"akari/models/user"

	"github.com/kataras/golog"
	"github.com/kataras/iris/v12"
)

func Signup(c iris.Context, route models.Route) {

	golog.Debug(route)

	if c.Method() != "POST" || route.Tertiary != "" || len(route.Tail) != 0 {
		c.StopWithStatus(iris.StatusNotFound)
		return
	}

	var signupForm SignupForm
	if err := c.ReadBody(&signupForm); err != nil || signupForm.Username == "" || signupForm.Email == "" || signupForm.Password == "" {
		c.StopWithJSON(iris.StatusBadRequest, iris.Map{"message": "Please fully fill in the signup form"})
		return
	}

	if err := auth.ValidUsername(signupForm.Username); err != "" {
		c.StopWithJSON(iris.StatusBadRequest, iris.Map{"message": err})
		return
	}

	if !user.CheckUsernameAvailability(signupForm.Username) {
		c.StopWithJSON(iris.StatusConflict, iris.Map{"message": "Username not available"})
		return
	}

	if !auth.ValidEmail(signupForm.Email) {
		c.StopWithJSON(iris.StatusBadRequest, iris.Map{"message": "Invalid Email"})
		return
	}

	if !user.CheckEmailAvailability(signupForm.Email) {
		c.StopWithJSON(iris.StatusConflict, iris.Map{"message": "Email not available"})
	}

	if err := auth.ValidPassword(signupForm.Password); err != "" {
		c.StopWithJSON(iris.StatusBadRequest, iris.Map{"message": err})
		return
	}

	var id int64
	if id = user.CreateAccount(signupForm.Email, signupForm.Username, signupForm.Password); id == -1 {
		c.StopWithStatus(iris.StatusInternalServerError)
		return
	}

	c.StopWithStatus(iris.StatusCreated)

}
