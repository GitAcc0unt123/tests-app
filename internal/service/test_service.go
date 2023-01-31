package service

import (
	"tests_app/internal/models"
	"tests_app/internal/repository"
)

type TestService struct {
	repo repository.Test
}

func NewTestService(repo repository.Test) *TestService {
	return &TestService{repo: repo}
}

func (t *TestService) Create(userId int, test models.Test) (int, error) {
	return t.repo.Create(userId, test)
}

func (t *TestService) GetAll() ([]models.TestResponse, error) {
	return t.repo.GetAll()
}

func (t *TestService) GetById(testId int) (models.TestResponse, error) {
	return t.repo.GetById(testId)
}

func (t *TestService) Delete(userId, testId int) error {
	return t.repo.Delete(userId, testId)
}

func (t *TestService) Update(userId, testId int, input models.UpdateTestInput) error {
	return t.repo.Update(userId, testId, input)
}
