package services

import (
	"delivery-service/models"
	"delivery-service/repository"
	"errors"
	"fmt"
	"os"
	"time"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("неверный email или пароль")
	ErrPasswordMismatch   = errors.New("пароли не совпадают")
	ErrUserExists         = errors.New("пользователь с таким email уже существует")
	ErrInvalidToken       = errors.New("недействительный токен")
	ErrMissingJWTSecret   = errors.New("отсутствует JWT_SECRET в переменных окружения")
)

const (
	bcryptCost     = 12
	tokenExpiresIn = time.Hour * 24 * 7
)

type AuthService struct {
	userRepo *repository.UserRepository
}

func NewAuthService(userRepo *repository.UserRepository) *AuthService {
	if userRepo == nil {
		panic("user repository is required")
	}
	return &AuthService{userRepo: userRepo}
}

func (s *AuthService) Register(req *models.RegisterRequest) (*models.AuthResponse, error) {
	if err := s.validateRegistration(req); err != nil {
		return nil, err
	}

	existingUser, err := s.userRepo.GetUserByEmail(req.Email)
	if err != nil && !errors.Is(err, repository.ErrUserNotFound) {
		return nil, fmt.Errorf("ошибка при проверке существующего пользователя: %w", err)
	}
	if existingUser != nil {
		return nil, ErrUserExists
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcryptCost)
	if err != nil {
		return nil, fmt.Errorf("ошибка при хешировании пароля: %w", err)
	}

	user := &models.User{
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
	}

	if err := s.userRepo.CreateUser(user); err != nil {
		return nil, fmt.Errorf("ошибка при создании пользователя: %w", err)
	}

	token, err := s.generateToken(user)
	if err != nil {
		return nil, fmt.Errorf("ошибка при генерации токена: %w", err)
	}

	return &models.AuthResponse{
		Token: token,
		User:  user,
	}, nil
}

func (s *AuthService) Login(req *models.LoginRequest) (*models.AuthResponse, error) {
	if err := s.validateLogin(req); err != nil {
		return nil, err
	}

	user, err := s.userRepo.GetUserByEmail(req.Email)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, fmt.Errorf("ошибка при поиске пользователя: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	token, err := s.generateToken(user)
	if err != nil {
		return nil, fmt.Errorf("ошибка при генерации токена: %w", err)
	}

	return &models.AuthResponse{
		Token: token,
		User:  user,
	}, nil
}

func (s *AuthService) ValidateToken(tokenString string) (*models.User, error) {
	if tokenString == "" {
		return nil, ErrInvalidToken
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("неожиданный метод подписи: %v", token.Header["alg"])
		}
		secret := os.Getenv("JWT_SECRET")
		if secret == "" {
			return nil, ErrMissingJWTSecret
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("ошибка при проверке токена: %w", err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if !s.validateTokenClaims(claims) {
			return nil, ErrInvalidToken
		}

		userID := int64(claims["user_id"].(float64))
		user, err := s.userRepo.GetUserByID(userID)
		if err != nil {
			return nil, fmt.Errorf("ошибка при получении пользователя: %w", err)
		}

		return user, nil
	}

	return nil, ErrInvalidToken
}

func (s *AuthService) UpdateUser(userID int64, req *models.UpdateUserRequest) error {
	if err := s.validateUpdateRequest(req); err != nil {
		return err
	}

	user, err := s.userRepo.UpdateUser(userID, req)
	if err != nil {
		return fmt.Errorf("ошибка при обновлении пользователя: %w", err)
	}

	if user == nil {
		return repository.ErrUserNotFound
	}

	return nil
}

func (s *AuthService) GetUserByID(userID int64) (*models.User, error) {
	if userID <= 0 {
		return nil, repository.ErrInvalidInput
	}

	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении пользователя: %w", err)
	}

	return user, nil
}

func (s *AuthService) validateRegistration(req *models.RegisterRequest) error {
	if req == nil {
		return repository.ErrInvalidInput
	}
	if req.Name == "" || req.Email == "" || req.Password == "" || req.ConfirmPassword == "" {
		return repository.ErrInvalidInput
	}
	if req.Password != req.ConfirmPassword {
		return ErrPasswordMismatch
	}
	if len(req.Password) < 8 {
		return errors.New("пароль должен содержать минимум 8 символов")
	}
	return nil
}

func (s *AuthService) validateLogin(req *models.LoginRequest) error {
	if req == nil || req.Email == "" || req.Password == "" {
		return repository.ErrInvalidInput
	}
	return nil
}

func (s *AuthService) validateUpdateRequest(req *models.UpdateUserRequest) error {
	if req == nil {
		return repository.ErrInvalidInput
	}
	// добавить проверки здесь
	return nil
}

func (s *AuthService) validateTokenClaims(claims jwt.MapClaims) bool {
	exp, ok := claims["exp"].(float64)
	if !ok {
		return false
	}
	return time.Now().Unix() < int64(exp)
}

func (s *AuthService) generateToken(user *models.User) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return "", ErrMissingJWTSecret
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"exp":     time.Now().Add(tokenExpiresIn).Unix(),
	})

	return token.SignedString([]byte(secret))
}
