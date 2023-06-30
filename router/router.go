package router

import (
	//user defined package
	"echo/handler"

	//third party package
	"github.com/labstack/echo/v4"
)

func Router() {
	e := echo.New()

	// Router
	e.POST("/signup", handler.Signup)
	e.POST("/login", handler.Login)
	
	// Start the server
	e.Start(":8000")
}
