//package http
//
//import (
//	"bytes"
//	"encoding/json"
//	"net/http"
//	"net/http/httptest"
//	"testing"
//
//	"github.com/labstack/echo/v4"
//	"github.com/stretchr/testify/assert"
//)
//
//func TestUserHandler_Signup_InvalidEmail(t *testing.T) {
//	// Setup
//	e := echo.New()
//	handler := NewUserHandler(nil) // We'll test validation without actual repository
//
//	// Test data with invalid email
//	requestBody := map[string]interface{}{
//		"email":    "invalid-email",
//		"password": "password123",
//		"role":     "analyst",
//	}
//	jsonBody, _ := json.Marshal(requestBody)
//
//	// Create request
//	req := httptest.NewRequest(http.MethodPost, "/signup", bytes.NewBuffer(jsonBody))
//	req.Header.Set("Content-Type", "application/json")
//	rec := httptest.NewRecorder()
//	c := e.NewContext(req, rec)
//
//	// Execute
//	err := handler.Signup(c)
//
//	// Assertions
//	assert.NoError(t, err)
//	assert.Equal(t, http.StatusBadRequest, rec.Code)
//
//	var response map[string]interface{}
//	err = json.Unmarshal(rec.Body.Bytes(), &response)
//	assert.NoError(t, err)
//	assert.Equal(t, "invalid email format", response["error"])
//}
//
//func TestUserHandler_Signup_InvalidPassword(t *testing.T) {
//	// Setup
//	e := echo.New()
//	handler := NewUserHandler(nil) // We'll test validation without actual repository
//
//	// Test data with invalid password (too short)
//	requestBody := map[string]interface{}{
//		"email":    "test@example.com",
//		"password": "12345",
//		"role":     "analyst",
//	}
//	jsonBody, _ := json.Marshal(requestBody)
//
//	// Create request
//	req := httptest.NewRequest(http.MethodPost, "/signup", bytes.NewBuffer(jsonBody))
//	req.Header.Set("Content-Type", "application/json")
//	rec := httptest.NewRecorder()
//	c := e.NewContext(req, rec)
//
//	// Execute
//	err := handler.Signup(c)
//
//	// Assertions
//	assert.NoError(t, err)
//	assert.Equal(t, http.StatusBadRequest, rec.Code)
//
//	var response map[string]interface{}
//	err = json.Unmarshal(rec.Body.Bytes(), &response)
//	assert.NoError(t, err)
//	assert.Equal(t, "password must be 6-60 characters", response["error"])
//}
//
//func TestUserHandler_Signup_MissingFields(t *testing.T) {
//	// Setup
//	e := echo.New()
//	handler := NewUserHandler(nil) // We'll test validation without actual repository
//
//	// Test data with missing email
//	requestBody := map[string]interface{}{
//		"password": "password123",
//		"role":     "analyst",
//	}
//	jsonBody, _ := json.Marshal(requestBody)
//
//	// Create request
//	req := httptest.NewRequest(http.MethodPost, "/signup", bytes.NewBuffer(jsonBody))
//	req.Header.Set("Content-Type", "application/json")
//	rec := httptest.NewRecorder()
//	c := e.NewContext(req, rec)
//
//	// Execute
//	err := handler.Signup(c)
//
//	// Assertions
//	assert.NoError(t, err)
//	assert.Equal(t, http.StatusBadRequest, rec.Code)
//
//	var response map[string]interface{}
//	err = json.Unmarshal(rec.Body.Bytes(), &response)
//	assert.NoError(t, err)
//	assert.Equal(t, "email and password are required", response["error"])
//}
//
//func TestUserHandler_Login_InvalidRequest(t *testing.T) {
//	// Setup
//	e := echo.New()
//	handler := NewUserHandler(nil) // We'll test validation without actual repository
//
//	// Create request with invalid JSON
//	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBufferString("invalid json"))
//	req.Header.Set("Content-Type", "application/json")
//	rec := httptest.NewRecorder()
//	c := e.NewContext(req, rec)
//
//	// Execute
//	err := handler.Login(c)
//
//	// Assertions
//	assert.NoError(t, err)
//	assert.Equal(t, http.StatusBadRequest, rec.Code)
//
//	var response map[string]interface{}
//	err = json.Unmarshal(rec.Body.Bytes(), &response)
//	assert.NoError(t, err)
//	assert.Equal(t, "invalid request payload", response["error"])
//}

package http

import (
	"bytes"
	"errors"
	"microservice-b/model"
	"microservice-b/utils"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// MockUserRepo mocks the UserRepository
type MockUserRepo struct {
	mock.Mock
}

func (m *MockUserRepo) Signup(u *model.SignupRequest) error {
	args := m.Called(u)
	return args.Error(0)
}

func (m *MockUserRepo) Login(email, password string) (string, error) {
	args := m.Called(email, password)
	return args.String(0), args.Error(1)
}

func TestUserHandler_Signup(t *testing.T) {
	e := echo.New()
	mockRepo := new(MockUserRepo)
	handler := NewUserHandler(mockRepo)

	tests := []struct {
		name           string
		payload        string
		mockReturn     error
		expectedStatus int
	}{
		{
			name:           "invalid bind",
			payload:        "{invalid json",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "missing email",
			payload:        `{"password":"password123"}`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "invalid email",
			payload:        `{"email":"invalid","password":"password123"}`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "invalid password",
			payload:        `{"email":"test@example.com","password":"123"}`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "name too long",
			payload:        `{"email":"test@example.com","password":"password123","first_name":"` + strings.Repeat("a", 256) + `"}`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "invalid role",
			payload:        `{"email":"test@example.com","password":"password123","role":"superuser"}`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "user already exists",
			payload:        `{"email":"exists@example.com","password":"password123"}`,
			mockReturn:     errors.New("user already exists"),
			expectedStatus: http.StatusConflict,
		},
		{
			name:           "internal server error",
			payload:        `{"email":"error@example.com","password":"password123"}`,
			mockReturn:     errors.New("db down"),
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:           "success signup",
			payload:        `{"email":"new@example.com","password":"password123"}`,
			mockReturn:     nil,
			expectedStatus: http.StatusCreated,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/signup", bytes.NewBufferString(tt.payload))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			// Reset expectations
			mockRepo.ExpectedCalls = nil

			// Only set up mock if this test expects a Signup call
			if tt.expectedStatus == http.StatusCreated || tt.expectedStatus == http.StatusConflict || tt.expectedStatus == http.StatusInternalServerError {
				mockRepo.On("Signup", mock.AnythingOfType("*model.SignupRequest")).Return(tt.mockReturn)
			}

			err := handler.Signup(c)
			require.NoError(t, err)
			require.Equal(t, tt.expectedStatus, rec.Code)
		})
	}
}

func TestUserHandler_Login(t *testing.T) {
	e := echo.New()
	mockRepo := new(MockUserRepo)
	handler := NewUserHandler(mockRepo)

	tests := []struct {
		name           string
		payload        string
		mockReturn     string
		mockErr        error
		expectedStatus int
	}{
		{
			name:           "invalid bind",
			payload:        "{invalid json",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "email not found",
			payload:        `{"email":"notfound@example.com","password":"pass"}`,
			mockErr:        utils.ErrEmailNotFound,
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "password mismatch",
			payload:        `{"email":"user@example.com","password":"wrongpass"}`,
			mockErr:        utils.ErrPasswordMismatch,
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "internal error",
			payload:        `{"email":"error@example.com","password":"pass"}`,
			mockErr:        errors.New("db down"),
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:           "success login",
			payload:        `{"email":"user@example.com","password":"pass"}`,
			mockReturn:     "token123",
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBufferString(tt.payload))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			// Setup mock
			mockRepo.ExpectedCalls = nil
			if tt.mockErr != nil || tt.mockReturn != "" {
				mockRepo.On("Login", mock.Anything, mock.Anything).Return(tt.mockReturn, tt.mockErr)
			}

			err := handler.Login(c)
			require.NoError(t, err)
			require.Equal(t, tt.expectedStatus, rec.Code)
		})
	}
}
