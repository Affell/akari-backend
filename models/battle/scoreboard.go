package battle

import (
	"akari/models/postgresql"
	"database/sql"

	"github.com/jackc/pgx/v4"
	"github.com/kataras/golog"
)

func GetScoreboard(offset int) (users []map[string]interface{}) {

	sqlCo, err := pgx.ConnectConfig(postgresql.SQLCtx, postgresql.SQLConn)
	if err != nil {
		return
	}

	defer sqlCo.Close(postgresql.SQLCtx)

	query := "SELECT username, score FROM account WHERE enable=true ORDER BY score DESC LIMIT 50 OFFSET $1"

	rows, err := sqlCo.Query(postgresql.SQLCtx, query, offset)
	if err != nil {
		golog.Errorf("execution query '%s':\n%s", query, err.Error())
		return
	}
	defer rows.Close()

	for rows.Next() {
		var (
			score    sql.NullInt64
			username sql.NullString
		)

		err = rows.Scan(
			&username,
			&score,
		)

		if err == pgx.ErrNoRows {
			return
		} else if err != nil {
			golog.Errorf("psql request '%v' failed with error : %v", query, err)
			return
		}

		users = append(users, map[string]interface{}{
			"username": username.String,
			"score":    score.Int64,
		})
	}

	return
}
