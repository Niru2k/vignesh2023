package authentication

import (
	//inbuilt package
	"errors"
	"net/http"
	"time"

	//third party package
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
)

// setting authentication and authorization for admin and user
func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		tokenString := c.Request().Header.Get("Authorization")
		if tokenString == "" {
			return c.String(http.StatusUnauthorized, "Missing token")
		}
		for index, char := range tokenString {
			if char == ' ' {
				tokenString = tokenString[index+1:]
			}
		}
		claims := jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte("secret"), nil // Replace "secret" with your own secret key
		})

		if err != nil || !token.Valid {
			return c.String(http.StatusUnauthorized, "Invalid token")
		}
		//Check whether token is expired or not
		a, ok := claims["exp"].(int64)
		if ok && a < time.Now().Unix() {
			return c.String(http.StatusUnauthorized, "Expired token")
		}

		c.Set("email", claims["email"])
		c.Set("role", claims["role"])

		return next(c)
	}
}

// Admin verifying authentication API
func AdminAuth(c echo.Context) error {
	role := c.Get("role").(string)
	if role != "admin" {
		return errors.New("only admins can access to this endpoint")
	}
	return nil
}

// User verifying authentication API
func UserAuth(c echo.Context) error {
	role := c.Get("role").(string)
	if role != "user" {
		return errors.New("only user can access to this endpoint")
	}
	return nil
}
