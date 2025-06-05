package usecase

import (
	"fmt"
	"github.com/fire9900/auth/internal/models"
	"github.com/fire9900/auth/pkg/auth"
)

func (uc *UserUseCase) Authenticate(email string, password string) (string, string, int64, error) {
	user, err := uc.repo.GetByEmail(email)
	if err != nil {
		if err == models.ErrorUserNotFound {
			return "", "", 0, models.ErrorWrongPassword
		}
		return "", "", 0, err
	}

	if err := user.CheckPassword(password); err != nil {
		return "", "", 0, models.ErrorWrongPassword
	}

	accessToken, expiresIn, err := auth.GenerateAccessToken(user.ID)
	if err != nil {
		return "", "", 0, fmt.Errorf("ошибка генерации access токена: %w", err)
	}

	refreshToken, err := auth.GenerateRefreshToken(user.ID)
	if err != nil {
		return "", "", 0, fmt.Errorf("ошибка генерации refresh токена: %w", err)
	}

	return accessToken, refreshToken, expiresIn, nil
}
