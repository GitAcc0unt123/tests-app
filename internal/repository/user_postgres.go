package repository

import (
	"fmt"
	"log"
	"strings"
	"tests_app/internal/models"

	"github.com/jmoiron/sqlx"
)

type AuthPostgres struct {
	db *sqlx.DB
}

func NewAuthPostgres(db *sqlx.DB) *AuthPostgres {
	return &AuthPostgres{db: db}
}

func (a *AuthPostgres) Create(user models.User) (int, error) {
	row := a.db.QueryRow(`
	INSERT INTO users (name, username, password, email, created_at)
	VALUES ($1, $2, $3, $4, now()) RETURNING id`,
		user.Name,
		user.Username,
		user.Password,
		user.Email)
	var id int
	if err := row.Scan(&id); err != nil {
		log.Print(err.Error())
		return 0, err
	}
	return id, nil
}

func (a *AuthPostgres) Get(username, password string) (models.User, error) {
	var user models.User
	err := a.db.Get(&user, `SELECT * FROM users WHERE username = $1 AND password = $2`, username, password)
	return user, err
}

func (a *AuthPostgres) Update(userId int, input models.UpdateUserInput) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argIndex := 1

	if input.Name != nil {
		setValues = append(setValues, fmt.Sprintf("name=$%d", argIndex))
		args = append(args, *input.Name)
		argIndex++
	}

	if input.Email != nil {
		setValues = append(setValues, fmt.Sprintf("email=$%d", argIndex))
		args = append(args, *input.Email)
		argIndex++
	}

	if input.Username != nil {
		setValues = append(setValues, fmt.Sprintf("username=$%d", argIndex))
		args = append(args, *input.Username)
		argIndex++
	}

	if input.Password != nil {
		setValues = append(setValues, fmt.Sprintf("password=$%d", argIndex))
		args = append(args, *input.Password)
		argIndex++
	}

	args = append(args, userId)
	setQuery := strings.Join(setValues, ", ")
	query := fmt.Sprintf("UPDATE users SET %s WHERE id=$%d", setQuery, argIndex)
	_, err := a.db.Exec(query, args...)
	return err
}
