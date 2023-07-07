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
	e.POST("/posting",handler.Jobposting)
	e.GET("/posting",handler.GetJobPostingDetails)
	e.GET("posting/:id",handler.GetJobPostingByID)
	e.PUT("posting/:id",handler.UpdateJob)
	e.DELETE("posting/:id",handler.DeleteJob)
	e.Logger.Fatal(e.Start(":8000"))
}
