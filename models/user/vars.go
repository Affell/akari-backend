package user

const PERMISSIONS_EDIT_PERMISSION = "edit.user.permission"

var SensibleFields []string = []string{
	"email",
	"password",
}

type User struct {
	ID       int64  `structs:"id"`
	Email    string `structs:"email"`
	Username string `structs:"username"`
	Password string `structs:"password"`
	Score    int64  `structs:"score"`
	Enable   bool   `structs:"enable"`
}
