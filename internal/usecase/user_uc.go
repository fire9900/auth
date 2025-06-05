package usecase

import (
	"fmt"
	"github.com/fire9900/auth/internal/models"
	"github.com/fire9900/auth/internal/repository"
)

type UseCase interface {
	GetAllUsers() ([]models.User, error)
	GetUserByID(id int) (models.User, error)
	GetUserByEmail(email string) (models.User, error)
	CreateUser(user models.User) (models.User, error)
	UpdateUser(id int, user models.User) (models.User, error)
	DeleteUser(id int) error
	CheckPassword(id int, password string) bool
	Authenticate(email string, password string) (string, string, int64, error)
}

type UserUseCase struct {
	repo repository.UserRepository
}

func NewUserUseCase(repo repository.UserRepository) *UserUseCase {
	return &UserUseCase{repo: repo}
}

func (uc *UserUseCase) GetAllUsers() ([]models.User, error) {
	return uc.repo.GetAll()
}

func (uc *UserUseCase) GetUserByID(id int) (models.User, error) {
	return uc.repo.GetByID(id)
}

func (uc *UserUseCase) GetUserByEmail(email string) (models.User, error) {
	return uc.repo.GetByEmail(email)
}

func (uc *UserUseCase) CreateUser(user models.User) (models.User, error) {
	if user.Password == "" {
		return models.User{}, fmt.Errorf("пароль не может быть пустым")
	}

	if err := user.HashPassword(); err != nil {
		return models.User{}, err
	}
	user.Role = "user"

	return uc.repo.Create(user)
}

func (uc *UserUseCase) UpdateUser(id int, user models.User) (models.User, error) {
	return uc.repo.Update(id, user)
}

func (uc *UserUseCase) DeleteUser(id int) error {
	return uc.repo.Delete(id)
}

func (uc *UserUseCase) CheckPassword(id int, password string) bool {
	return uc.repo.CheckPassword(id, password)
}
