package service

import (
	"tests_app/internal/models"
	"tests_app/internal/repository"
)

type TestAnswerService struct {
	repo repository.TestAnswer
}

func NewTestAnswerService(repo repository.TestAnswer) *TestAnswerService {
	return &TestAnswerService{repo: repo}
}

func (t *TestAnswerService) Create(testAnswer models.TestAnswer) error {
	return t.repo.Create(testAnswer)
}

func (t *TestAnswerService) Get(userId, testId int) (models.TestAnswer, error) {
	return t.repo.Get(userId, testId)
}
