package repository

import (
	//third party package

	//user defined package
	// "echo/handler"
	"echo/helper"
	"echo/models"
)

func CreateTables() {
	helper.Db.AutoMigrate(&models.Information{})
	helper.Db.AutoMigrate(&models.Jobposting{})
	helper.Db.AutoMigrate(&models.Comments{})
}
