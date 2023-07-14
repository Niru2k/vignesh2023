package main


import (
	//user defined package
	"echo/router"
       "echo/driver"
	   "echo/repository"
	)

func main() {
	driver.DatabaseConnection()
	repository.CreateTables()
	router.Router()
}


