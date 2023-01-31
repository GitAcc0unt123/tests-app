package service

import (
	"tests_app/internal/models"
	"tests_app/internal/repository"
)

type QuestionAnswerService struct {
	repo repository.QuestionAnswer
}

func NewQuestionAnswerService(repo repository.QuestionAnswer) *QuestionAnswerService {
	return &QuestionAnswerService{repo: repo}
}

func (t *QuestionAnswerService) Upsert(userId int, questionAnswer models.UpsertQuestionAnswerInput) error {
	return t.repo.Upsert(userId, questionAnswer)
}

func (t *QuestionAnswerService) GetAnswerByQuestionId(userId, questionId int) (models.QuestionAnswer, error) {
	return t.repo.GetAnswerByQuestionId(userId, questionId)
}

func (t *QuestionAnswerService) GetAnswersByTestId(userId, testId int) ([]models.QuestionAnswer, error) {
	return t.repo.GetAnswersByTestId(userId, testId)
}
