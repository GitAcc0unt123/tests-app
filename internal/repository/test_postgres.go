package repository

import (
	"fmt"
	"strings"
	"tests_app/internal/models"

	"github.com/jmoiron/sqlx"
)

type TestPostgres struct {
	db *sqlx.DB
}

func NewTestPostgres(db *sqlx.DB) *TestPostgres {
	return &TestPostgres{db: db}
}

func (t *TestPostgres) Create(userId int, test models.Test) (int, error) {
	row := t.db.QueryRow(`
	INSERT INTO tests (title, description, random_questions_order, questions_visibility, start_time, end_time, duration_sec)
	VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`,
		test.Title,
		test.Description,
		test.RandomQuestionsOrder,
		test.QuestionsVisibility,
		test.Start,
		test.End,
		test.Duration_sec)
	var id int
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func (t *TestPostgres) GetAll() ([]models.TestResponse, error) {
	var tests []models.TestResponse
	err := t.db.Select(&tests, `SELECT * FROM tests`)
	if err != nil {
		return nil, err
	}
	return tests, nil
}

func (t *TestPostgres) GetById(testId int) (models.TestResponse, error) {
	var test models.TestResponse
	err := t.db.Get(&test, `SELECT * FROM tests WHERE id = $1`, testId)
	return test, err
}

func (t *TestPostgres) Delete(userId, testId int) error {
	_, err := t.db.Exec(`DELETE FROM tests WHERE id = $1`, testId)
	return err
}

func (t *TestPostgres) Update(userId, testId int, input models.UpdateTestInput) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argIndex := 1

	if input.Title != nil {
		setValues = append(setValues, fmt.Sprintf("title=$%d", argIndex))
		args = append(args, *input.Title)
		argIndex++
	}

	if input.Description != nil {
		setValues = append(setValues, fmt.Sprintf("description=$%d", argIndex))
		args = append(args, *input.Description)
		argIndex++
	}

	if input.RandomQuestionsOrder != nil {
		setValues = append(setValues, fmt.Sprintf("random_questions_order=$%d", argIndex))
		args = append(args, *input.RandomQuestionsOrder)
		argIndex++
	}

	if input.QuestionsVisibility != nil {
		setValues = append(setValues, fmt.Sprintf("questions_visibility=$%d", argIndex))
		args = append(args, *input.QuestionsVisibility)
		argIndex++
	}

	if input.Start != nil {
		setValues = append(setValues, fmt.Sprintf("start_time=%d", argIndex))
		args = append(args, *input.Start)
		argIndex++
	}

	if input.End != nil {
		setValues = append(setValues, fmt.Sprintf("end_time=%d", argIndex))
		args = append(args, *input.End)
		argIndex++
	}

	if input.Duration_sec != nil {
		setValues = append(setValues, fmt.Sprintf("duration_sec=%d", argIndex))
		args = append(args, *input.Duration_sec)
		argIndex++
	}

	args = append(args, testId)
	setQuery := strings.Join(setValues, ", ")
	query := fmt.Sprintf("UPDATE tests SET %s WHERE id=$%d", setQuery, argIndex)
	_, err := t.db.Exec(query, args...)
	return err
}
