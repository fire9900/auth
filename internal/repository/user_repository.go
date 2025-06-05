package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/fire9900/auth/internal/models"
	"github.com/fire9900/auth/pkg/logger"
	"go.uber.org/zap"
)

type UserRepository interface {
	GetAll() ([]models.User, error)
	GetByID(id int) (models.User, error)
	GetByEmail(email string) (models.User, error)
	Create(user models.User) (models.User, error)
	Update(id int, user models.User) (models.User, error)
	Delete(id int) error
	CheckPassword(id int, password string) bool
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) GetAll() ([]models.User, error) {
	logger.Logger.Info("Получение всех пользователей")

	query := `SELECT id, name, email, password, role FROM users ORDER BY id`
	rows, err := r.db.Query(query)
	if err != nil {
		logger.Logger.Error("Ошибка при запросе всех пользователей",
			zap.Error(err),
			zap.String("метод", "GetAll"))
		return nil, fmt.Errorf("ошибка при запросе всех пользователей: %w", err)
	}
	defer func() {
		if err := rows.Close(); err != nil {
			logger.Logger.Warn("Ошибка при закрытии rows",
				zap.Error(err),
				zap.String("метод", "GetAll"))
		}
	}()

	var users []models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Role)
		if err != nil {
			logger.Logger.Error("Ошибка при сканировании данных пользователя",
				zap.Error(err),
				zap.String("метод", "GetAll"))
			return nil, fmt.Errorf("ошибка при сканировании данных пользователя: %w", err)
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		logger.Logger.Error("Ошибка при итерации по результатам запроса",
			zap.Error(err),
			zap.String("метод", "GetAll"))
		return nil, fmt.Errorf("ошибка при итерации по результатам запроса: %w", err)
	}

	logger.Logger.Info("Успешно получены все пользователи",
		zap.Int("количество", len(users)))
	return users, nil
}

func (r *userRepository) GetByID(id int) (models.User, error) {
	logger.Logger.Info("Получение пользователя по ID",
		zap.Int("id", id))

	query := `SELECT id, name, email, password, role FROM users WHERE id = $1`
	row := r.db.QueryRow(query, id)

	var user models.User
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Role)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.Logger.Warn("Пользователь не найден",
				zap.Int("id", id),
				zap.String("метод", "GetByID"))
			return models.User{}, models.ErrorUserNotFound
		}
		logger.Logger.Error("Ошибка при получении пользователя по ID",
			zap.Error(err),
			zap.Int("id", id),
			zap.String("метод", "GetByID"))
		return models.User{}, fmt.Errorf("ошибка при получении пользователя по ID: %w", err)
	}

	logger.Logger.Info("Успешно получен пользователь по ID",
		zap.Int("id", id))
	return user, nil
}

func (r *userRepository) Create(user models.User) (models.User, error) {
	logger.Logger.Info("Создание нового пользователя",
		zap.String("email", user.Email))

	query := `INSERT INTO users (name, email, password) 
		 VALUES ($1, $2, $3)
		 RETURNING id, name, email, password`

	var createdUser models.User
	err := r.db.QueryRow(
		query,
		user.Name,
		user.Email,
		user.Password,
	).Scan(
		&createdUser.ID,
		&createdUser.Name,
		&createdUser.Email,
		&createdUser.Password,
	)

	if err != nil {
		logger.Logger.Error("Ошибка при создании пользователя",
			zap.Error(err),
			zap.String("email", user.Email),
			zap.String("метод", "Create"))
		return models.User{}, fmt.Errorf("ошибка при создании пользователя: %w", err)
	}

	logger.Logger.Info("Пользователь успешно создан",
		zap.Int("id", createdUser.ID),
		zap.String("email", createdUser.Email))
	return createdUser, nil
}

func (r *userRepository) GetByEmail(email string) (models.User, error) {
	query := `SELECT id, name, email, password, role FROM users WHERE email = $1`
	row := r.db.QueryRow(query, email)

	var user models.User
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Role)
	if err != nil {
		fmt.Println(err.Error())
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, models.ErrorUserNotFound
		}
		return models.User{}, fmt.Errorf("Ошибка получения пользователя по Email: %w", err)
	}
	fmt.Println(user)
	return user, nil
}

func (r *userRepository) Update(id int, user models.User) (models.User, error) {
	logger.Logger.Info("Обновление данных пользователя",
		zap.Int("id", id),
		zap.String("email", user.Email))

	query := `UPDATE users
		 SET name = $1, email = $2, password = $3
		 WHERE id = $4
		 RETURNING id, name, email, password`

	var updatedUser models.User
	err := r.db.QueryRow(
		query,
		user.Name,
		user.Email,
		user.Password,
		id,
	).Scan(
		&updatedUser.ID,
		&updatedUser.Name,
		&updatedUser.Email,
		&updatedUser.Password,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.Logger.Warn("Пользователь не найден при обновлении",
				zap.Int("id", id),
				zap.String("метод", "Update"))
			return models.User{}, models.ErrorUserNotFound
		}
		logger.Logger.Error("Ошибка при обновлении пользователя",
			zap.Error(err),
			zap.Int("id", id),
			zap.String("метод", "Update"))
		return models.User{}, fmt.Errorf("ошибка при обновлении пользователя: %w", err)
	}

	if err = updatedUser.HashPassword(); err != nil {
		logger.Logger.Error("Ошибка при хешировании пароля",
			zap.Error(err),
			zap.Int("id", id),
			zap.String("метод", "Update"))
		return models.User{}, err
	}

	logger.Logger.Info("Данные пользователя успешно обновлены",
		zap.Int("id", updatedUser.ID),
		zap.String("email", updatedUser.Email))
	return updatedUser, nil
}

func (r *userRepository) Delete(id int) error {
	logger.Logger.Info("Удаление пользователя",
		zap.Int("id", id))

	query := `DELETE FROM users WHERE id = $1`
	result, err := r.db.Exec(query, id)
	if err != nil {
		logger.Logger.Error("Ошибка при удалении пользователя",
			zap.Error(err),
			zap.Int("id", id),
			zap.String("метод", "Delete"))
		return fmt.Errorf("ошибка при удалении пользователя: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		logger.Logger.Error("Ошибка при получении количества удаленных строк",
			zap.Error(err),
			zap.Int("id", id),
			zap.String("метод", "Delete"))
		return fmt.Errorf("ошибка при получении количества удаленных строк: %w", err)
	}

	if rowsAffected == 0 {
		logger.Logger.Warn("Пользователь не найден при удалении",
			zap.Int("id", id),
			zap.String("метод", "Delete"))
		return models.ErrorUserNotFound
	}

	logger.Logger.Info("Пользователь успешно удален",
		zap.Int("id", id),
		zap.Int64("rowsAffected", rowsAffected))
	return nil
}

func (r *userRepository) CheckPassword(id int, password string) bool {
	logger.Logger.Debug("Проверка пароля пользователя",
		zap.Int("id", id))

	user, err := r.GetByID(id)
	if err != nil {
		logger.Logger.Error("Ошибка при проверке пароля - пользователь не найден",
			zap.Error(err),
			zap.Int("id", id),
			zap.String("метод", "CheckPassword"))
		return false
	}

	isValid := user.VerifyPassword(password)
	if !isValid {
		logger.Logger.Warn("Неверный пароль пользователя",
			zap.Int("id", id),
			zap.String("метод", "CheckPassword"))
	} else {
		logger.Logger.Debug("Пароль пользователя проверен успешно",
			zap.Int("id", id))
	}

	return isValid
}
