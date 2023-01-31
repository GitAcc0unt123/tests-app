package repository

import (
	"fmt"
	"log"
	"strings"
	"tests_app/internal/models"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type QuestionPostgres struct {
	db *sqlx.DB
}

func NewQuestionPostgres(db *sqlx.DB) *QuestionPostgres {
	return &QuestionPostgres{db: db}
}

func (t *QuestionPostgres) Create(userId int, Question models.Question) (int, error) {
	row := t.db.QueryRow(`
	INSERT INTO questions (test_id, text, answer_type, show_answers, true_answers)
	SELECT $1, $2, $3, $4, $5
	WHERE now() < (select start_time from tests where id = $1) RETURNING id`,
		Question.TestId,
		Question.Text,
		Question.AnswerType,
		pq.Array(Question.ShowAnswers),
		pq.Array(Question.TrueAnswers))

	var id int
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func (t *QuestionPostgres) GetAll(userId, testId int) ([]models.GetQuestionResponse, error) {
	type tempGetQuestionResponse struct {
		Id          int            `db:"id"`
		TestId      int            `db:"test_id"`
		Text        string         `db:"text"`
		AnswerType  string         `db:"answer_type"`
		ShowAnswers pq.StringArray `db:"show_answers"`
	}

	var output []tempGetQuestionResponse
	err := t.db.Select(&output, `SELECT id, test_id, text, answer_type, show_answers FROM questions WHERE test_id = $1 ORDER BY id`, testId)
	if err != nil {
		return nil, err
	}

	questions := make([]models.GetQuestionResponse, len(output))
	for i, row := range output {
		questions[i].Id = row.Id
		questions[i].TestId = row.TestId
		questions[i].Text = row.Text
		questions[i].AnswerType = row.AnswerType
		questions[i].ShowAnswers = row.ShowAnswers
	}

	return questions, nil
}

func (t *QuestionPostgres) GetAllWithAnswer(userId, testId int) ([]models.GetQuestionResponse2, error) {
	type tempGetQuestionResponse struct {
		Id          int            `db:"id"`
		TestId      int            `db:"test_id"`
		Text        string         `db:"text"`
		AnswerType  string         `db:"answer_type"`
		ShowAnswers pq.StringArray `db:"show_answers"`
		TrueAnswers pq.StringArray `db:"true_answers"`
	}

	var output []tempGetQuestionResponse
	err := t.db.Select(&output, `SELECT * FROM questions WHERE test_id = $1 ORDER BY id`, testId)
	if err != nil {
		return nil, err
	}

	questions := make([]models.GetQuestionResponse2, len(output))
	for i, row := range output {
		questions[i].Id = row.Id
		questions[i].TestId = row.TestId
		questions[i].Text = row.Text
		questions[i].AnswerType = row.AnswerType
		questions[i].ShowAnswers = row.ShowAnswers
		questions[i].TrueAnswers = row.TrueAnswers
	}

	return questions, nil
}

func (t *QuestionPostgres) GetById(userId, questionId int) (models.Question, error) {
	var q models.Question
	query := `SELECT * FROM questions WHERE id = $1`
	// wrap the output parameter in pq.Array for receiving into it
	err := t.db.QueryRow(query, questionId).Scan(
		&q.Id,
		&q.TestId,
		&q.Text,
		&q.AnswerType,
		pq.Array(&q.ShowAnswers),
		pq.Array(&q.TrueAnswers))
	if err != nil {
		log.Fatal(err)
	}
	//err := t.db.Get(&question, `SELECT * FROM questions WHERE id=$1`, questionId)
	return q, err
}

// Первый вопрос без ответа в заданном порядке
func (t *QuestionPostgres) GetNext(userId, testId int) (models.GetQuestionResponse, error) {
	var question models.GetQuestionResponse

	err := t.db.Get(&question, `
	SELECT * FROM questions
	WHERE test_id = $1 AND id NOT IN (SELECT question_id FROM question_answers WHERE user_id = $2 AND test_id = $1)
	ORDER BY id
	LIMIT 1`,
		testId,
		userId)
	return question, err
}

func (t *QuestionPostgres) Delete(userId, questionId int) error {
	_, err := t.db.Exec(`DELETE FROM questions WHERE id = $1`, questionId)
	return err
}

func (t *QuestionPostgres) Update(userId, questionId int, input models.UpdateQuestionInput) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argIndex := 1

	if input.Text != nil {
		setValues = append(setValues, fmt.Sprintf("text=$%d", argIndex))
		args = append(args, *input.Text)
		argIndex++
	}

	if input.AnswerType != nil {
		setValues = append(setValues, fmt.Sprintf("answer_type=$%d", argIndex))
		args = append(args, *input.AnswerType)
		argIndex++
	}

	if input.ShowAnswers != nil {
		showAnswers := pq.Array(*input.ShowAnswers)
		setValues = append(setValues, fmt.Sprintf("show_answers=$%d", argIndex))
		args = append(args, showAnswers)
		argIndex++
	}

	if input.TrueAnswers != nil {
		trueAnswers := pq.Array(*input.TrueAnswers)
		setValues = append(setValues, fmt.Sprintf("true_answers=$%d", argIndex))
		args = append(args, trueAnswers)
		argIndex++
	}

	args = append(args, questionId)
	setQuery := strings.Join(setValues, ", ")
	query := fmt.Sprintf("UPDATE questions SET %s WHERE id=$%d", setQuery, argIndex)
	_, err := t.db.Exec(query, args...)
	return err
}
