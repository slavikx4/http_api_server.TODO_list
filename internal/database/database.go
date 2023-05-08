package database

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/slavikx4/http-api-server/internal/config"
	"log"
	"time"
)

var DB *pgxpool.Pool

type User struct {
	Login string
	Tasks []string
}

func init() {
	var err error
	DB, err = pgxpool.New(context.Background(), fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		config.Config.PgUser, config.Config.PgPassword, config.Config.PgHost, config.Config.PgPort, config.Config.PgBase))
	if err != nil {
		log.Fatalln(err)
	}
	if err := DB.Ping(context.Background()); err != nil {
		log.Fatalln(err)
	}
}

func AddUserInBase(inputLogin, inputPassword string) error {

	query := `INSERT INTO "User"("Login", "Password") VALUES ($1,$2)`

	if _, err := DB.Exec(context.Background(), query, inputLogin, inputPassword); err != nil {
		return err
	}

	return nil
}

func CheckUserInBase(inputLogin string) (string, error) {

	var password string

	query := `SELECT "Password" FROM "User" WHERE "Login"=$1`

	row := DB.QueryRow(context.Background(), query, inputLogin)

	if err := row.Scan(&password); err != nil {
		return "", err
	}

	return password, nil
}

func AddTaskInBase(login, task string) error {

	query := `INSERT INTO "Task" ("User", "Text","Data", "Stage") VALUES($1,$2,$3, $4)`

	if _, err := DB.Exec(context.Background(), query, login, task, time.Now(), "waiting"); err != nil {
		return err
	}

	return nil
}

func UnloadTaskFromBase(login string) (*User, error) {

	var user = User{
		Login: login,
		Tasks: make([]string, 0),
	}

	query := `SELECT ("Text") FROM "Task" WHERE "User"=$1`

	rows, err := DB.Query(context.Background(), query, login)
	if err != nil {
		return nil, err
	}

	var text string

	for rows.Next() {

		if err := rows.Scan(&text); err != nil {
			return nil, err
		}
		user.Tasks = append(user.Tasks, text)
	}

	return &user, nil
}
