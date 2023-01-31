package repository

import (
	"log"
	"tests_app/internal/models"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type QuestionAnswerPostgres struct {
	db *sqlx.DB
}

func NewQuestionAnswerPostgres(db *sqlx.DB) *QuestionAnswerPostgres {
	return &QuestionAnswerPostgres{db: db}
}

func (t *QuestionAnswerPostgres) Upsert(userId int, questionAnswer models.UpsertQuestionAnswerInput) error {
	row := t.db.QueryRow(`
		INSERT INTO question_answers (user_id, question_id, answer, time)
		VALUES($1, $2, $3, now())
		ON CONFLICT (user_id, question_id)
		DO
			UPDATE SET answer = $3, time = now()`,
		userId,
		questionAnswer.QuestionId,
		pq.Array(questionAnswer.Answer))
	return row.Err()
}

func (t *QuestionAnswerPostgres) GetAnswerByQuestionId(userId, questionId int) (models.QuestionAnswer, error) {
	query := `
	SELECT * FROM question_answers
	WHERE user_id = $1 AND question_id = $2`

	var a models.QuestionAnswer
	// wrap the output parameter in pq.Array for receiving into it
	err := t.db.QueryRow(query, userId, questionId).Scan(&a.UserId, &a.QuestionId, pq.Array(&a.Answer), &a.Time)
	if err != nil {
		log.Fatal(err)
	}

	return a, err
}

func (t *QuestionAnswerPostgres) GetAnswersByTestId(userId, testId int) ([]models.QuestionAnswer, error) {
	query := `
	SELECT * FROM question_answers
	WHERE user_id = $1 AND question_id IN (SELECT id FROM questions WHERE test_id = $2)`

	type TempQuestionAnswer struct {
		UserId     int            `db:"user_id"`
		QuestionId int            `db:"question_id"`
		Answer     pq.StringArray `db:"answer"`
		Time       time.Time      `db:"time"`
	}

	var output []TempQuestionAnswer
	err := t.db.Select(&output, query, userId, testId)
	if err != nil {
		return nil, err
	}

	answers := make([]models.QuestionAnswer, len(output))
	for i, row := range output {
		answers[i].UserId = row.UserId
		answers[i].QuestionId = row.QuestionId
		answers[i].Answer = row.Answer
		answers[i].Time = row.Time
	}

	return answers, err
}
