package service

import (
	"tests_app/internal/models"
	"tests_app/internal/repository"

	"github.com/google/uuid"
)

type RefreshSessionService struct {
	repo repository.RefreshSession
}

func NewRefreshSessionService(repo repository.RefreshSession) *RefreshSessionService {
	return &RefreshSessionService{repo: repo}
}

func (r *RefreshSessionService) Create(refreshSession models.RefreshSession) error {
	return r.repo.Create(refreshSession)
}

func (r *RefreshSessionService) Get(refreshToken uuid.UUID) (models.RefreshSession, error) {
	return r.repo.Get(refreshToken)
}

func (r *RefreshSessionService) Revoke(refreshToken uuid.UUID) error {
	return r.repo.Revoke(refreshToken)
}
