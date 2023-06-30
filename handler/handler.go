package handler

import (
	"fmt"
	"net/http"

	"echo/models"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
)

func jwtMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		tokenString := c.Request().Header.Get("Authorization")
		if tokenString == "" {
			return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Missing token"})
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Invalid token")
			}
			return models.SigningKey, nil
		})
		if err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Invalid token"})
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			userID := claims["id"].(string)
			for _, user := range models.Users {
				if user.ID == userID {
					// Store user information in the context
					c.Set("user", user)
					return next(c)
				}
			}
		}

		return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Invalid token"})
	}
}

func Signup(c echo.Context) error {
	var user models.User
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request payload"})
	}

	// Check if the user already exists by ID, Email, or Username
	for _, u := range models.Users {
		if u.ID == user.ID || u.Email == user.Email || u.Username == user.Username {
			return c.JSON(http.StatusConflict, map[string]string{"message": "User already exists"})
		}
	}

	// Generate JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":    user.ID,
		"email": user.Email,
	})

	tokenString, err := token.SignedString(models.SigningKey)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to generate token"})
	}

	// Save user to the list
	models.Users = append(models.Users, user)
	fmt.Println("store", models.Users)
	return c.JSON(http.StatusOK, map[string]string{"token": tokenString})
}

func Login(c echo.Context) error {
	var user models.User
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request payload"})
	}

	// Find user by username and validate if the user already exists
	for _, u := range models.Users {
		if u.Username != user.Username {
			return c.JSON(http.StatusUnauthorized, map[string]string{"message": "user not found"})
		}
		if u.Password != user.Password {
			return c.JSON(http.StatusUnauthorized, map[string]string{"message": "incorrect password"})
		}

	}
	return c.JSON(http.StatusAccepted, map[string]string{"message": "login successful"})
}

func Welcome(c echo.Context) error {
	return c.String(http.StatusOK, "signup and login")
}
