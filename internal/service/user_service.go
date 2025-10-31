package service

import (
"crypto/sha256"
"fmt"

"github.com/ethandiaz/yad-back/internal/domain"
)

type UserService struct {
	repo domain.UserRepository
}

func NewUserService(repo domain.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) Register(email, firstName, lastName, password string) (*domain.User, error) {
	existing, _ := s.repo.GetByEmail(email)
	if existing != nil {
		return nil, fmt.Errorf("email already registered")
	}

	hashedPassword := s.hashPassword(password)
	user := domain.NewUser(email, firstName, lastName, hashedPassword)

	if err := s.repo.Create(user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

func (s *UserService) GetUser(id string) (*domain.User, error) {
	return s.repo.GetByID(id)
}

func (s *UserService) GetUserByEmail(email string) (*domain.User, error) {
	return s.repo.GetByEmail(email)
}

func (s *UserService) UpdateUser(user *domain.User) error {
	return s.repo.Update(user)
}

func (s *UserService) DeleteUser(id string) error {
	return s.repo.Delete(id)
}

func (s *UserService) hashPassword(password string) string {
	hash := sha256.Sum256([]byte(password))
	return fmt.Sprintf("%x", hash)
}

func (s *UserService) VerifyPassword(user *domain.User, password string) bool {
	return user.Password == s.hashPassword(password)
}
