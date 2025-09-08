package usecase

import (
	"fmt"
	"microservice-b/internal/repository"
	"microservice-b/middleware"
	"microservice-b/model"
	"microservice-b/utils"
	"regexp"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type IUserRepository interface {
	Signup(u *model.SignupRequest) error
	Login(email, password string) (string, error)
}

type UserRepository struct {
	Repo      *repository.UserRepo
	JWTSecret string
}

// Signup
func (s *UserRepository) Signup(u *model.SignupRequest) error {
	// Check if user already exists
	existingUser, _ := s.Repo.GetByEmail(u.Email)
	if existingUser != nil {
		return fmt.Errorf("user with email %s already exists", u.Email)
	}

	// Hash password
	hashed, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashed)

	// Save to DB
	return s.Repo.CreateUser(u)
}

// Login
func (s *UserRepository) Login(email, password string) (string, error) {
	u, err := s.Repo.GetByEmail(email)
	if err != nil {
		if err == utils.ErrEmailNotFound {
			return "", utils.ErrEmailNotFound
		}
		return "", err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); err != nil {
		return "", utils.ErrPasswordMismatch
	}

	// Generate JWT token
	token, err := middleware.GenerateJWT(u.ID, u.Email, u.Role, s.JWTSecret, 24)
	if err != nil {
		return "", err
	}

	// Update last_login
	_ = s.Repo.UpdateLastLogin(u.ID)

	return token, nil
}

// IsValidEmail returns true if the email has a valid format
func IsValidEmail(email string) bool {
	email = strings.TrimSpace(strings.ToLower(email))
	emailRegex := `^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`
	matched, _ := regexp.MatchString(emailRegex, email)
	return matched
}

func IsValidPassword(password string) bool {
	length := len(password)
	return length >= 6 && length <= 60
}
