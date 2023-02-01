package usecase

import (
	"context"

	"github.com/VetKA-org/vetka/internal/repo"
)

type AuthUseCase struct {
	usersRepo repo.Users
}

func NewAuthUseCase(users repo.Users) *AuthUseCase {
	return &AuthUseCase{users}
}

func (uc *AuthUseCase) Login(ctx context.Context, username, password string) error {
	return nil
}
