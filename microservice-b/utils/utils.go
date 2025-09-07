package utils

import (
	"errors"
	"microservice-b/model"

	"github.com/labstack/echo/v4"
)

var (
	ErrEmailNotFound    = errors.New("email is wrong")
	ErrPasswordMismatch = errors.New("password is mismatch")
)

// ErrorResponse sends a structured error response using the provided status code, message, optional code, and details.
func ErrorResponse(ctx echo.Context, statusCode int, errMsg string, errCode int, details string) error {
	return ctx.JSON(statusCode, model.ErrorResponse{
		Error:   errMsg,
		Code:    errCode,
		Details: details,
	})
}
