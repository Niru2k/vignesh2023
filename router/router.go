package router

import (
	"echo/handler"

	"github.com/labstack/echo/v4"
)

func Router() {
	e := echo.New()

	// Routes
	e.POST("/signup", handler.Signup)
	e.POST("/login", handler.Login)
    e.GET("/",handler.Welcome)
	// Start the server
	e.Start(":8000")
}
