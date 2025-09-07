package http

import (
	"microservice-b/internal/usecase"
	"microservice-b/model"
	"microservice-b/utils"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	userRepo *usecase.UserRepository
}

func NewUserHandler(repo *usecase.UserRepository) *UserHandler {
	return &UserHandler{userRepo: repo}
}

// Signup godoc
// @Summary Create a new user account
// @Description Register a new user with email, password, optional name fields, and role. The email must be unique. Passwords must be between 6 and 60 characters. Role defaults to "analyst" if not provided.
// @Tags Users
// @Accept json
// @Produce json
// @Param user body model.SignupRequest true "User signup payload"
// @Success 201 {object} model.SignupResponse "User created successfully"
// @Failure 400 {object} model.ErrorResponse "Invalid request or validation failed"
// @Failure 409 {object} model.ErrorResponse "User with email already exists"
// @Failure 500 {object} model.ErrorResponse "Internal server error"
// @Router /signup [post]
func (c *UserHandler) Signup(ctx echo.Context) error {
	u := new(model.SignupRequest)
	if err := ctx.Bind(u); err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "invalid request payload", 1001, "")
	}

	// Trim spaces
	u.FirstName = strings.TrimSpace(u.FirstName)
	u.LastName = strings.TrimSpace(u.LastName)
	u.Email = strings.TrimSpace(u.Email)
	u.Password = strings.TrimSpace(u.Password)
	u.Role = strings.TrimSpace(u.Role)

	// Validate required fields
	if u.Email == "" || u.Password == "" {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "email and password are required", 1002, "")
	}

	// Validate email and password
	if !usecase.IsValidEmail(u.Email) {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "invalid email format", 1003, "")
	}

	if !usecase.IsValidPassword(u.Password) {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "password must be 6-60 characters", 1004, "")
	}

	// Validate optional names (max 255 chars)
	if len(u.FirstName) > 255 || len(u.LastName) > 255 {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "first_name/last_name too long", 1005, "")
	}

	// Validate role
	if u.Role == "" {
		u.Role = "analyst"
	} else if u.Role != "admin" && u.Role != "analyst" {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "role must be 'admin' or 'analyst'", 1006, "")
	}

	// Call service to create user
	if err := c.userRepo.Signup(u); err != nil {
		if strings.Contains(err.Error(), "already exists") {
			return utils.ErrorResponse(ctx, http.StatusConflict, err.Error(), 1007, "")
		}
		return utils.ErrorResponse(ctx, http.StatusInternalServerError, "internal server error", 1008, err.Error())
	}

	return ctx.JSON(http.StatusCreated, echo.Map{"message": "user created"})
}

// Login godoc
// @Summary Authenticate a user and return a JWT token
// @Description Authenticates a user by validating the provided email and password. Returns a JWT token upon successful login.
// @Tags Users
// @Accept json
// @Produce json
// @Param credentials body model.Login true "User login credentials payload"
// @Success 200 {object}  model.LoginResponse "JWT token response"
// @Failure 400 {object} model.ErrorResponse "invalid request payload"}
// @Failure 401 {object} model.ErrorResponse "invalid credentials"}
// @Failure 500 {object} model.ErrorResponse "something went wrong"}
// @Router /login [post]
func (c *UserHandler) Login(ctx echo.Context) error {
	req := model.Login{}
	if err := ctx.Bind(&req); err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "invalid request payload", 2001, "")
	}
	token, err := c.userRepo.Login(req.Email, req.Password)
	if err != nil {
		switch err {
		case utils.ErrEmailNotFound:
			return utils.ErrorResponse(ctx, http.StatusUnauthorized, "email is wrong", 2002, "")
		case utils.ErrPasswordMismatch:
			return utils.ErrorResponse(ctx, http.StatusUnauthorized, "password is mismatch", 2003, "")
		default:
			return utils.ErrorResponse(ctx, http.StatusInternalServerError, "something went wrong", 2004, err.Error())
		}
	}
	return ctx.JSON(http.StatusOK, model.LoginResponse{Token: token})
}
