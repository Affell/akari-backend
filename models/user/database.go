package user

import (
	"akari/models/postgresql"
	"database/sql"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/kataras/golog"
)

func GetSQLUserToken(username, password string) (token UserToken, err error) {

	var (
		id    sql.NullInt64
		email sql.NullString
	)

	query := "select id, email from account " +
		"where enable=true and username=$1 and password=crypt($2, password)"

	sqlCo, err := pgx.ConnectConfig(postgresql.SQLCtx, postgresql.SQLConn)
	if err != nil {
		return
	}
	defer sqlCo.Close(postgresql.SQLCtx)

	err = sqlCo.QueryRow(postgresql.SQLCtx, query, username, password).Scan(
		&id,
		&email,
	)

	if err == pgx.ErrNoRows {
		golog.Infof("Tentative de connexion infructueuse pour l'utilisateur : %v. Email inconnu ou utilisateur non actif.", username)
		return
	} else if err != nil {
		golog.Error(query, err)
		return
	}

	token = UserToken{
		ID:        id.Int64,
		Username:  username,
		Email:     email.String,
		CreatedAt: time.Now(),
	}

	return
}

func CreateAccount(email, username, password string) int64 {

	query := "insert into account (email, username, password, enable) " +
		"VALUES ($1,$2,crypt($3, gen_salt('bf')),true) RETURNING id"

	sqlCo, err := pgx.ConnectConfig(postgresql.SQLCtx, postgresql.SQLConn)
	if err != nil {
		golog.Error(err)
		return -1
	}
	defer sqlCo.Close(postgresql.SQLCtx)

	row := sqlCo.QueryRow(postgresql.SQLCtx, query, email, username, password)
	var id sql.NullInt64
	err = row.Scan(&id)
	if err != nil {
		golog.Error(err)
		return -1
	}
	return id.Int64
}

func DeleteAccount(id int64) (msg string) {

	query := "update account set enable=false where id=$1"

	sqlCo, err := pgx.ConnectConfig(postgresql.SQLCtx, postgresql.SQLConn)
	if err != nil {
		msg = "Internal server error"
		return
	}
	defer sqlCo.Close(postgresql.SQLCtx)

	cmd, err := sqlCo.Exec(postgresql.SQLCtx, query, id)
	if err != nil {
		return "Internal server error"
	} else if cmd.RowsAffected() == 0 {
		return "Account not found"
	}
	return
}

func SecurityCheck(id int64, password string) (checked bool) {
	query := "select id from account where enable=true and id=$1 and password=crypt($2, password)"

	sqlCo, err := pgx.ConnectConfig(postgresql.SQLCtx, postgresql.SQLConn)
	if err != nil {
		return
	}
	defer sqlCo.Close(postgresql.SQLCtx)

	var id_ int64
	err = sqlCo.QueryRow(postgresql.SQLCtx, query, id, password).Scan(&id_)
	if err == nil {
		return id_ == id
	}
	return
}

func CheckUsernameAvailability(username string) (available bool) {
	query := "select id from account where username=$1"

	sqlCo, err := pgx.ConnectConfig(postgresql.SQLCtx, postgresql.SQLConn)
	if err != nil {
		return
	}
	defer sqlCo.Close(postgresql.SQLCtx)

	var id int64
	err = sqlCo.QueryRow(postgresql.SQLCtx, query, username).Scan(&id)
	if err == pgx.ErrNoRows {
		return true
	}
	return
}

func CheckEmailAvailability(email string) (available bool) {
	query := "select id from account where email=$1"

	sqlCo, err := pgx.ConnectConfig(postgresql.SQLCtx, postgresql.SQLConn)
	if err != nil {
		return
	}
	defer sqlCo.Close(postgresql.SQLCtx)

	var id int64
	err = sqlCo.QueryRow(postgresql.SQLCtx, query, email).Scan(&id)
	if err == pgx.ErrNoRows {
		return true
	}
	return
}

func GetUserById(UserId int64) (u User, err error) {

	var (
		email, username, password sql.NullString
	)

	sqlCo, err := pgx.ConnectConfig(postgresql.SQLCtx, postgresql.SQLConn)
	if err != nil {
		return
	}
	defer sqlCo.Close(postgresql.SQLCtx)

	// Requête GetUserById
	var query = "SELECT email, username, password " +
		"FROM account " +
		"WHERE enable=true and id=$1 "

	err = sqlCo.QueryRow(postgresql.SQLCtx, query, UserId).Scan(
		&email,
		&username,
		&password,
	)

	if err != nil {
		return
	}

	u = User{
		ID:       UserId,
		Email:    email.String,
		Username: username.String,
		Password: password.String,
	}

	return
}

func GetAllUsers() (users []User) {

	sqlCo, err := pgx.ConnectConfig(postgresql.SQLCtx, postgresql.SQLConn)
	if err != nil {
		return
	}

	defer sqlCo.Close(postgresql.SQLCtx)

	query := "SELECT id, email, username, password, enable FROM account"

	rows, err := sqlCo.Query(postgresql.SQLCtx, query)
	if err != nil {
		golog.Errorf("execution query '%s':\n%s", query, err.Error())
		return
	}
	defer rows.Close()

	for rows.Next() {
		var (
			id                        sql.NullInt64
			email, username, password sql.NullString
			enable                    sql.NullBool
		)

		err = rows.Scan(
			&id,
			&email,
			&username,
			&password,
			&enable,
		)

		if err == pgx.ErrNoRows {
			return
		} else if err != nil {
			golog.Errorf("psql request '%v' failed with error : %v", query, err)
			return
		}

		users = append(users, User{
			ID:       id.Int64,
			Email:    email.String,
			Username: username.String,
			Password: password.String,
			Enable:   enable.Bool,
		})
	}

	return
}

func GetUserByEmail(userEmail string) (u User, msg string) {

	if userEmail == "" {
		msg = "Empty email"
		return
	}

	sqlCo, err := pgx.ConnectConfig(postgresql.SQLCtx, postgresql.SQLConn)
	if err != nil {
		msg = "Internal server error"
		return
	}
	defer sqlCo.Close(postgresql.SQLCtx)

	query := "SELECT id, username, password " +
		"FROM account " +
		"where email=$1"

	var (
		id                 sql.NullInt64
		username, password sql.NullString
	)
	err = sqlCo.QueryRow(postgresql.SQLCtx, query, userEmail).Scan(
		&id,
		&username,
		&password,
	)

	if err == pgx.ErrNoRows {
		msg = "Username not found"
		return
	} else if err != nil {
		golog.Errorf("psql request '%v' failed with error : %v", query, err)
		msg = "Internal server error"
		return
	}

	u = User{
		ID:       id.Int64,
		Email:    userEmail,
		Username: username.String,
		Password: password.String,
	}

	return
}

func UpdateUser(user User, password bool) (ok bool) {

	sqlCo, err := pgx.ConnectConfig(postgresql.SQLCtx, postgresql.SQLConn)
	if err != nil {
		return
	}
	defer sqlCo.Close(postgresql.SQLCtx)

	var query string
	var args []interface{}
	if password {
		query = "UPDATE account set (email, username, password) = ($1,$2,crypt($3, gen_salt('bf'))) " +
			"WHERE id=$4"
		args = []interface{}{
			user.Email,
			user.Username,
			user.Password,
			user.ID,
		}
	} else {
		query = "UPDATE account set (username, email) = ($1,$2) " +
			"WHERE id=$3"
		args = []interface{}{
			user.Username,
			user.Email,
			user.ID,
		}
	}

	cmd, err := sqlCo.Exec(postgresql.SQLCtx, query, args...)
	ok = cmd.RowsAffected() == 1 && err == nil
	return
}
