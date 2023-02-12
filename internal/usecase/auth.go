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
	rolesRepo repo.Roles
}

func NewAuthUseCase(secret string, users repo.Users, roles repo.Roles) *AuthUseCase {
	return &AuthUseCase{secret, users, roles}
}

func (uc *AuthUseCase) Login(ctx context.Context, login, password string) (entity.JWTToken, error) {
	user, err := uc.usersRepo.Verify(ctx, login, password)
	if err != nil {
		if errors.Is(err, entity.ErrInvalidCredentials) {
			return "", err
		}

		return "", fmt.Errorf("AuthUseCase - Login - uc.usersRepo.Verify: %w", err)
	}

	roles, err := uc.rolesRepo.Get(ctx, user.ID)
	if err != nil {
		return "", fmt.Errorf("AuthUseCase - Login - uc.rolesRepo.Get: %w", err)
	}

	token, err := entity.NewJWTToken(user, roles, uc.secret, _defaultTokenLifeTime)
	if err != nil {
		return "", fmt.Errorf("AuthUseCase - Login - entity.NewJWTToken: %w", err)
	}

	return token, nil
}
