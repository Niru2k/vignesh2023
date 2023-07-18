package router

import (
	//user defined package
	"echo/authentication"
	"echo/handler"

	//third party package
	"github.com/labstack/echo/v4"
)

func Router() {
	e := echo.New()
	//SignUp Api
	e.POST("/signup", handler.Signup)  
	//Login Api                                                              
	e.POST("/login", handler.Login)     
	//only admin has authorization                                                                         
	e.POST("/jobposting", handler.Jobposting, authentication.AuthMiddleware)    
	e.PUT("updatejobpostbyid/:id", handler.UpdateJob, authentication.AuthMiddleware)     
	e.DELETE("deletejobpostbyid/:id", handler.DeleteJob, authentication.AuthMiddleware) 
	//both admin and user has the access                      
	e.GET("/getjobposts", handler.GetJobPostingDetails, authentication.AuthMiddleware)           
	e.GET("getjobpostbyid/:id", handler.GetJobPostingByID, authentication.AuthMiddleware)                       
	e.GET("getjobpostbycompanyname/:companyname", handler.GetJobPostingByCompany, authentication.AuthMiddleware) 
	e.GET("/getallcomments", handler.GetUserComments, authentication.AuthMiddleware)                          
	e.GET("/getcommentsbyid/:id", handler.GetCommentByID, authentication.AuthMiddleware)                         
	//only user has authorization
	e.POST("/user/comments", handler.UserComments, authentication.AuthMiddleware)                 
	e.PUT("updatecommentbyid/:id", handler.UpdateComment, authentication.AuthMiddleware)                       
	e.DELETE("deletecommentbyid/:id", handler.DeleteCommentById, authentication.AuthMiddleware)   
    e.Logger.Fatal(e.Start(":8000"))
}
