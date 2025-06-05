package handlers

import (
	"github.com/fire9900/auth/internal/models"
	"github.com/fire9900/auth/internal/usecase"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// @title sigma Auth API
// @version 6.0
// @description API аутентификации пользователей
// @host localhost:8080
// @BasePath /api/v1
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

type UserHandler struct {
	userUseCase usecase.UseCase
}

func NewUserHandler(useCase usecase.UseCase) *UserHandler {
	return &UserHandler{userUseCase: useCase}
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// @Summary Получить всех пользователей
// @Description Получить список пользователей
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {array} models.User
// @Failure 400 {object} object
// @Router /users [get]
func (h *UserHandler) GetAll(c *gin.Context) {
	users, err := h.userUseCase.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, users)
}

func (h *UserHandler) Logout(c *gin.Context) {
	c.Set("userID", nil)
	c.JSON(http.StatusOK, gin.H{"details": "Выход из аккаунта"})
	c.Next()
}

// @Summary Найти пользователя по ID
// @Description В url запроса помещается ID пользователя, если он существует, возвращается объект пользователя
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "ID пользователя"
// @Success 200 {object} models.User
// @Failure 400 {object} object
// @Failure 404 {object} object
// @Router /users/{id} [get]
func (h *UserHandler) GetByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Параметр id недействительный"})
		return
	}

	user, err := h.userUseCase.GetUserByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}

// @Summary Получить пользователя по email
// @Description В запрос устанавливается email и получается пользователь, если он существует
// @Tags users
// @Accept json
// @Produce json
// @Param email path string true "Email пользователя"
// @Success 200 {object} models.User
// @Failure 400 {object} object
// @Failure 404 {object} object
// @Router /users/email/{email} [get]
func (h *UserHandler) GetByEmail(c *gin.Context) {
	email := c.Param("email")
	if email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Недействительный email"})
		return
	}

	user, err := h.userUseCase.GetUserByEmail(email)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Пользователь не найден"})
		return
	}
	c.JSON(http.StatusOK, user)
}

// @Summary Создать пользователя
// @Description Создать нового пользователя
// @Tags users
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param user body models.User true "Данные пользователя"
// @Success 201 {object} models.User
// @Failure 400 {object} object
// @Failure 401 {object} object
// @Failure 500 {object} object
// @Router /users [post]
func (h *UserHandler) Create(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if user.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Пароль не может быть пустым"})
		return
	}

	createUser, err := h.userUseCase.CreateUser(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, createUser)
}

// @Summary Обновить пароль пользователя
// @Description Обновляет пароль пользователя по ID
// @Tags users
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "ID пользователя"
// @Param password body object{password=string} true "Новый пароль"
// @Success 200 {object} models.User
// @Failure 400 {object} object
// @Failure 401 {object} object
// @Failure 404 {object} object
// @Failure 500 {object} object
// @Router /users/{id}/password [put]
func (h *UserHandler) UpdatePassword(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Параметр id недействительный"})
		return
	}

	var updateData struct {
		Password string `json:"password"`
	}

	if err = c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.userUseCase.GetUserByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	user.Password = updateData.Password
	if err = user.HashPassword(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ошибка обработки пароля"})
		return
	}
	updateUser, err := h.userUseCase.UpdateUser(id, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	updateUser.Password = ":)"

	c.JSON(http.StatusOK, updateUser)
}

// @Summary Удалить пользователя
// @Description Удаляет пользователя по ID
// @Tags users
// @Security ApiKeyAuth
// @Param id path int true "ID пользователя"
// @Success 200 {object} object{message=string}
// @Failure 400 {object} object
// @Failure 401 {object} object
// @Failure 500 {object} object
// @Router /users/{id} [delete]
func (h *UserHandler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.userUseCase.DeleteUser(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Пользователь с id: " + strconv.Itoa(id) + " удален"})
}

// @Summary Проверить пароль
// @Description Проверяет соответствие пароля пользователя
// @Tags users
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "ID пользователя"
// @Param password body object{password=string} true "Пароль для проверки"
// @Success 200 {boolean} bool
// @Failure 400 {object} object
// @Failure 401 {object} object
// @Router /users/{id}/check-password [post]
func (h *UserHandler) CheckPassword(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var InPassword struct {
		Password string `json:"password"`
	}

	if err = c.ShouldBindJSON(&InPassword); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	booler := h.userUseCase.CheckPassword(id, InPassword.Password)
	c.JSON(http.StatusOK, booler)
}
