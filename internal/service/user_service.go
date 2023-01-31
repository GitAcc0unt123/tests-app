package service

import (
	"tests_app/internal/models"
	"tests_app/internal/repository"
	"tests_app/pkg/hash"
)

type UserService struct {
	repo   repository.User
	hasher hash.PasswordHasher
}

func NewUserService(repo repository.User, hasher hash.PasswordHasher) *UserService {
	return &UserService{
		repo:   repo,
		hasher: hasher}
}

func (s *UserService) Get(username, password string) (models.User, error) {
	password, err := s.hasher.Hash(password)
	if err != nil {
		return models.User{}, err
	}
	return s.repo.Get(username, password)
}

func (s *UserService) Create(user models.User) (int, error) {
	password, err := s.hasher.Hash(user.Password)
	if err != nil {
		return 0, err
	}

	user.Password = password
	return s.repo.Create(user)
}

func (s *UserService) Update(userId int, input models.UpdateUserInput) error {
	if input.Password != nil {
		password, err := s.hasher.Hash(*input.Password) // ?
		if err != nil {
			return err
		}

		input.Password = &password
	}
	return s.repo.Update(userId, input)
}
