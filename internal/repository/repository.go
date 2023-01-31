package repository

import (
	"tests_app/internal/models"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type User interface {
	Create(user models.User) (int, error)
	Get(username, password string) (models.User, error)
	Update(userId int, input models.UpdateUserInput) error
}

type Test interface {
	Create(userId int, test models.Test) (int, error)
	GetAll() ([]models.TestResponse, error)
	GetById(testId int) (models.TestResponse, error)
	Delete(userId, testId int) error
	Update(userId, testId int, input models.UpdateTestInput) error
}

type Question interface {
	Create(userId int, question models.Question) (int, error)
	GetAll(userId, testId int) ([]models.GetQuestionResponse, error)
	GetAllWithAnswer(userId, testId int) ([]models.GetQuestionResponse2, error)
	GetById(userId, questionId int) (models.Question, error)
	GetNext(userId, testId int) (models.GetQuestionResponse, error)
	Delete(userId, questionId int) error
	Update(userId, questionId int, input models.UpdateQuestionInput) error
}

type QuestionAnswer interface {
	Upsert(userId int, questionAnswer models.UpsertQuestionAnswerInput) error
	GetAnswerByQuestionId(userId, questionId int) (models.QuestionAnswer, error)
	GetAnswersByTestId(userId, testId int) ([]models.QuestionAnswer, error)
}

type TestAnswer interface {
	Create(testAnswer models.TestAnswer) error
	Get(userId, testId int) (models.TestAnswer, error)
}

type RefreshSession interface {
	Create(refreshSessions models.RefreshSession) error
	Get(refreshToken uuid.UUID) (models.RefreshSession, error)
	Revoke(refreshToken uuid.UUID) error
}

type Repository struct {
	User
	Test
	Question
	QuestionAnswer
	TestAnswer
	RefreshSession
}

func New(db *sqlx.DB) *Repository {
	return &Repository{
		User:           NewAuthPostgres(db),
		Test:           NewTestPostgres(db),
		Question:       NewQuestionPostgres(db),
		QuestionAnswer: NewQuestionAnswerPostgres(db),
		TestAnswer:     NewTestAnswerPostgres(db),
		RefreshSession: NewRefreshSessionPostgres(db),
	}
}
