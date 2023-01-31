package service

import (
	"tests_app/internal/models"
	"tests_app/internal/repository"
)

type QuestionService struct {
	repo repository.Question
}

func NewQuestionService(repo repository.Question) *QuestionService {
	return &QuestionService{repo: repo}
}

func (t *QuestionService) Create(userId int, question models.Question) (int, error) {
	return t.repo.Create(userId, question)
}

func (t *QuestionService) GetAll(userId, testId int) ([]models.GetQuestionResponse, error) {
	return t.repo.GetAll(userId, testId)
}

func (t *QuestionService) GetAllWithAnswer(userId, testId int) ([]models.GetQuestionResponse2, error) {
	return t.repo.GetAllWithAnswer(userId, testId)
}

func (t *QuestionService) GetById(userId, questionId int) (models.Question, error) {
	return t.repo.GetById(userId, questionId)
}

func (t *QuestionService) GetNext(userId, questionId int) (models.GetQuestionResponse, error) {
	return t.repo.GetNext(userId, questionId)
}

func (t *QuestionService) Delete(userId, questionId int) error {
	return t.repo.Delete(userId, questionId)
}

func (t *QuestionService) Update(userId, questionId int, input models.UpdateQuestionInput) error {
	return t.repo.Update(userId, questionId, input)
}
