package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

// GenerateJWT generates a signed JWT token
func GenerateJWT(userID uint64, email, role, secret string, expiryHours int) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"email":   email,
		"role":    role,
		"exp":     time.Now().Add(time.Duration(expiryHours) * time.Hour).Unix(),
	}
	fmt.Println("secret---", secret)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// JWTMiddleware returns Echo JWT middleware configured with the secret
func JWTMiddleware(secret string) echo.MiddlewareFunc {
	return echojwt.WithConfig(echojwt.Config{
		SigningKey:  []byte(secret),
		TokenLookup: "header:Authorization:Bearer ",
		ErrorHandler: func(c echo.Context, err error) error {
			c.Logger().Errorf("JWT validation failed: %v", err)
			return c.JSON(http.StatusUnauthorized, echo.Map{
				"error": "unauthorized",
			})
		},
	})
}
