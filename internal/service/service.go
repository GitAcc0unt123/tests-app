package service

import (
	"tests_app/internal/models"
	"tests_app/internal/repository"
	"tests_app/pkg/hash"
	"tests_app/pkg/token"

	"github.com/google/uuid"
)

//go:generate mockgen -source=service.go -destination=mocks/mock.go

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
	GetAnswerByQuestionId(userId, questionId int) (models.QuestionAnswer, error)
	GetAnswersByTestId(userId, testId int) ([]models.QuestionAnswer, error)
	Upsert(userId int, input models.UpsertQuestionAnswerInput) error
}

type TestAnswer interface {
	Create(testAnswer models.TestAnswer) error
	Get(userId, testId int) (models.TestAnswer, error)
}

type RefreshSession interface {
	Create(refreshSession models.RefreshSession) error
	Get(refreshToken uuid.UUID) (models.RefreshSession, error)
	Revoke(refreshToken uuid.UUID) error
}

type Service struct {
	User
	Test
	Question
	QuestionAnswer
	TestAnswer
	RefreshSession
	TokenManager *token.Manager
}

func New(repos *repository.Repository, tokenManager *token.Manager, hasher hash.PasswordHasher) *Service {
	return &Service{
		User:           NewUserService(repos.User, hasher),
		Test:           NewTestService(repos.Test),
		Question:       NewQuestionService(repos.Question),
		QuestionAnswer: NewQuestionAnswerService(repos.QuestionAnswer),
		TestAnswer:     NewTestAnswerService(repos.TestAnswer),
		RefreshSession: NewRefreshSessionService(repos.RefreshSession),
		TokenManager:   tokenManager,
	}
}
