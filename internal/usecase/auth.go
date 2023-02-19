package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/VetKA-org/vetka/internal/repo"
	"github.com/VetKA-org/vetka/pkg/entity"
	"github.com/VetKA-org/vetka/pkg/logger"
)

const (
	_defaultMinimalSecretLength = 32
	_defaultTokenLifeTime       = 15 * time.Minute
)

type AuthUseCase struct {
	secret    entity.Secret
	usersRepo repo.Users
	rolesRepo repo.Roles
}

func NewAuthUseCase(
	log *logger.Logger,
	secret entity.Secret,
	users repo.Users,
	roles repo.Roles,
) *AuthUseCase {
	if len([]byte(secret)) < _defaultMinimalSecretLength {
		log.Warn().Msg("Insecure signature: secret key is shorter than 32 bytes!")
	}

	return &AuthUseCase{secret, users, roles}
}

func (uc *AuthUseCase) Login(ctx context.Context, login, password string) (entity.JWTToken, error) {
	user, err := uc.usersRepo.Verify(ctx, login, password)
	if err != nil {
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
