package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/VetKA-org/vetka/internal/entity"
	"github.com/VetKA-org/vetka/internal/repo"
)

const _defaultTokenLifeTime = 15 * time.Minute

type AuthUseCase struct {
	secret    string
	usersRepo repo.Users
}

func NewAuthUseCase(secret string, users repo.Users) *AuthUseCase {
	return &AuthUseCase{secret, users}
}

func (uc *AuthUseCase) Login(ctx context.Context, login, password string) (entity.JWTToken, error) {
	user, err := uc.usersRepo.Verify(ctx, login, password)
	if err != nil {
		if errors.Is(err, entity.ErrInvalidCredentials) {
			return "", err
		}

		return "", fmt.Errorf("AuthUseCase - Login - uc.usersRepo.Verify: %w", err)
	}

	token, err := entity.NewJWTToken(user, uc.secret, _defaultTokenLifeTime)
	if err != nil {
		return "", fmt.Errorf("AuthUseCase - Login - entity.NewJWTToken: %w", err)
	}

	return token, nil
}
