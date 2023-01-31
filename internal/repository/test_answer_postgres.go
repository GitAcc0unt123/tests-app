package repository

import (
	"tests_app/internal/models"

	"github.com/jmoiron/sqlx"
)

type TestAnswerPostgres struct {
	db *sqlx.DB
}

func NewTestAnswerPostgres(db *sqlx.DB) *TestAnswerPostgres {
	return &TestAnswerPostgres{db: db}
}

func (t *TestAnswerPostgres) Create(testAnswer models.TestAnswer) error {
	row := t.db.QueryRow(`
	INSERT INTO test_answers (user_id, test_id, complete_time)
	VALUES ($1, $2, $3)`,
		testAnswer.UserId,
		testAnswer.TestId,
		testAnswer.CompleteTime)

	return row.Err()
}

func (t *TestAnswerPostgres) Get(userId, testId int) (models.TestAnswer, error) {
	var testAnswer models.TestAnswer
	err := t.db.Get(&testAnswer, `SELECT * FROM test_answers WHERE user_id = $1 AND test_id = $2`, userId, testId)
	return testAnswer, err
}
