package handler

import (
	//built in package
	"fmt"
	"net/http"

	//user defined package
	"echo/models"

	//third party package
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
)

// for signing up
func Signup(c echo.Context) error {
	var user models.User
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request payload"})
	}
	// Check if the user already exists by ID, Email, or Username
	for _, u := range models.Users {
		fmt.Printf("%+v", u)
		fmt.Printf("%+v", user)
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

// for login
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
