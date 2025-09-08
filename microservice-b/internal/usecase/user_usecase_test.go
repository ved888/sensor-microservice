package usecase

import (
	"errors"
	"testing"

	"microservice-b/middleware"
	"microservice-b/model"
	"microservice-b/utils"

	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

// UserRepoInterface defines the interface for user repository
type UserRepoInterface interface {
	CreateUser(user *model.SignupRequest) error
	GetByEmail(email string) (*model.User, error)
	UpdateLastLogin(id uint64) error
}

// MockUserRepo is a mock implementation of UserRepoInterface
type MockUserRepo struct {
	mock.Mock
}

func (m *MockUserRepo) CreateUser(user *model.SignupRequest) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepo) GetByEmail(email string) (*model.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepo) UpdateLastLogin(id uint64) error {
	args := m.Called(id)
	return args.Error(0)
}

// TestUserRepository is a test version of UserRepository that uses the interface
type TestUserRepository struct {
	Repo      UserRepoInterface
	JWTSecret string
}

func (s *TestUserRepository) Signup(u *model.SignupRequest) error {
	// Validate email format
	if !IsValidEmail(u.Email) {
		return errors.New("invalid email format")
	}

	// Validate password length
	if !IsValidPassword(u.Password) {
		return errors.New("password must be 6-60 characters")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("failed to hash password")
	}

	// Update the user object with hashed password
	u.Password = string(hashedPassword)

	// Create user in repository
	return s.Repo.CreateUser(u)
}

func (s *TestUserRepository) Login(email, password string) (string, error) {
	// Get user by email
	user, err := s.Repo.GetByEmail(email)
	if err != nil {
		return "", err
	}

	// Check password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", utils.ErrPasswordMismatch
	}

	// Update last login
	err = s.Repo.UpdateLastLogin(user.ID)
	if err != nil {
		return "", err
	}

	// Generate JWT token
	token, err := middleware.GenerateJWT(user.ID, user.Email, user.Role, s.JWTSecret, 24)
	if err != nil {
		return "", err
	}

	return token, nil
}

func TestUserRepository_Signup_Success(t *testing.T) {
	// Setup
	mockRepo := new(MockUserRepo)
	userRepo := &TestUserRepository{
		Repo:      mockRepo,
		JWTSecret: "test-secret-key",
	}

	// Test data
	user := &model.SignupRequest{
		Email:    "test@example.com",
		Password: "password123",
		Role:     "analyst",
	}

	// Mock expectations
	mockRepo.On("CreateUser", mock.AnythingOfType("*model.SignupRequest")).Return(nil)

	// Execute
	err := userRepo.Signup(user)

	// Assertions
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestUserRepository_Signup_EmailAlreadyExists(t *testing.T) {
	// Setup
	mockRepo := new(MockUserRepo)
	userRepo := &TestUserRepository{
		Repo:      mockRepo,
		JWTSecret: "test-secret-key",
	}

	// Test data
	user := &model.SignupRequest{
		Email:    "existing@example.com",
		Password: "password123",
		Role:     "analyst",
	}

	// Mock expectations
	mockRepo.On("CreateUser", mock.AnythingOfType("*model.SignupRequest")).Return(errors.New("email already exists"))

	// Execute
	err := userRepo.Signup(user)

	// Assertions
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "email already exists")
	mockRepo.AssertExpectations(t)
}

func TestUserRepository_Login_Success(t *testing.T) {
	// Setup
	mockRepo := new(MockUserRepo)
	userRepo := &TestUserRepository{
		Repo:      mockRepo,
		JWTSecret: "test-secret-key",
	}

	// Test data
	email := "test@example.com"
	password := "password123"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	user := &model.User{
		ID:       1,
		Email:    email,
		Password: string(hashedPassword),
		Role:     "analyst",
	}

	// Mock expectations
	mockRepo.On("GetByEmail", email).Return(user, nil)
	mockRepo.On("UpdateLastLogin", uint64(1)).Return(nil)

	// Execute
	token, err := userRepo.Login(email, password)

	// Assertions
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	// Verify JWT token
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte("test-secret-key"), nil
	})
	assert.NoError(t, err)
	assert.True(t, parsedToken.Valid)

	claims := parsedToken.Claims.(jwt.MapClaims)
	assert.Equal(t, float64(1), claims["user_id"])
	assert.Equal(t, email, claims["email"])
	assert.Equal(t, "analyst", claims["role"])

	mockRepo.AssertExpectations(t)
}

func TestUserRepository_Login_EmailNotFound(t *testing.T) {
	// Setup
	mockRepo := new(MockUserRepo)
	userRepo := &TestUserRepository{
		Repo:      mockRepo,
		JWTSecret: "test-secret-key",
	}

	// Test data
	email := "nonexistent@example.com"
	password := "password123"

	// Mock expectations
	mockRepo.On("GetByEmail", email).Return(nil, utils.ErrEmailNotFound)

	// Execute
	token, err := userRepo.Login(email, password)

	// Assertions
	assert.Error(t, err)
	assert.Equal(t, utils.ErrEmailNotFound, err)
	assert.Empty(t, token)
	mockRepo.AssertExpectations(t)
}

func TestUserRepository_Login_WrongPassword(t *testing.T) {
	// Setup
	mockRepo := new(MockUserRepo)
	userRepo := &TestUserRepository{
		Repo:      mockRepo,
		JWTSecret: "test-secret-key",
	}

	// Test data
	email := "test@example.com"
	password := "wrongpassword"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("correctpassword"), bcrypt.DefaultCost)

	user := &model.User{
		ID:       1,
		Email:    email,
		Password: string(hashedPassword),
		Role:     "analyst",
	}

	// Mock expectations
	mockRepo.On("GetByEmail", email).Return(user, nil)

	// Execute
	token, err := userRepo.Login(email, password)

	// Assertions
	assert.Error(t, err)
	assert.Equal(t, utils.ErrPasswordMismatch, err)
	assert.Empty(t, token)
	mockRepo.AssertExpectations(t)
}

func TestIsValidEmail(t *testing.T) {
	tests := []struct {
		name     string
		email    string
		expected bool
	}{
		{"Valid email", "test@example.com", true},
		{"Valid email with subdomain", "user@mail.example.com", true},
		{"Invalid email - no @", "testexample.com", false},
		{"Invalid email - no domain", "test@", false},
		{"Invalid email - no local", "@example.com", false},
		{"Empty email", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsValidEmail(tt.email)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestIsValidPassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		expected bool
	}{
		{"Valid password", "password123", true},
		{"Valid password - minimum length", "123456", true},
		{"Valid password - maximum length", repeatString("a", 60), true},
		{"Invalid password - too short", "12345", false},
		{"Invalid password - too long", repeatString("a", 61), false},
		{"Empty password", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsValidPassword(tt.password)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// Helper function for string repetition
func repeatString(s string, count int) string {
	result := ""
	for i := 0; i < count; i++ {
		result += s
	}
	return result
}