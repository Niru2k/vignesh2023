package main

import (
	//user defined package
	"echo/driver"
	"echo/repository"
	"echo/router"
)

func main() {
	driver.DatabaseConnection()
	repository.CreateTables()
	router.Router()
}
