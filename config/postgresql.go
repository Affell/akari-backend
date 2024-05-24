package config

import (
	"context"
	"os"

	"github.com/jackc/pgx/v4"
	"github.com/kataras/golog"
)

func InitPgSQL() (context.Context, *pgx.ConnConfig) {

	connstring := "postgresql://"

	if env := os.Getenv("POSTGRES_USER"); env == "" {
		golog.Fatal("Bad 'POSTGRES_USER' parameter env")
	} else {
		connstring += env
	}

	if env := os.Getenv("POSTGRES_PASSWORD"); env == "" {
		golog.Warn("'POSTGRES_PASSWORD' not set parameter env")
	} else {
		connstring += ":" + env

	}

	if env := os.Getenv("POSTGRES_HOST"); env == "" {
		golog.Fatal("Bad 'POSTGRES_HOST' parameter env")
		os.Exit(1)
	} else {
		connstring += "@" + env
	}

	if env := os.Getenv("POSTGRES_DB"); env == "" {
		golog.Fatal("Bad 'POSTGRES_DB' parameter env")
	} else {
		connstring += "/" + env
	}

	connstring += "?sslmode=disable"

	connConf, err := pgx.ParseConfig(connstring)
	if err != nil {
		golog.Fatalf("Parse error : %s", err)
	}

	sqlCo, err := pgx.ConnectConfig(context.Background(), connConf)
	if err != nil {
		golog.Errorf("error connect psql : %s", err)
		return context.Background(), connConf
	}
	defer sqlCo.Close(context.Background())

	query := `
	CREATE EXTENSION IF NOT EXISTS pgcrypto;
	CREATE TABLE IF NOT EXISTS account (
		id 					SERIAL,
		email 			TEXT NOT NULL UNIQUE,
		username 		TEXT NOT NULL UNIQUE,
		password 		TEXT,
		enable			boolean,
		PRIMARY KEY(id)
	);
	CREATE TABLE IF NOT EXISTS score (
		account				BIGINT NOT NULL,
		score 			INT NOT NULL,
		PRIMARY KEY(account),
		FOREIGN KEY(account) REFERENCES account(id)
	);
	`

	_, err = sqlCo.Exec(context.Background(), query)
	if err != nil {
		golog.Errorf("%v : %v", query, err)
	}
	return context.Background(), connConf
}
