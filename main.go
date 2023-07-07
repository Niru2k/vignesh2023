package main

//user defined package
import (
	"echo/router"
       "echo/driver"
	   "echo/repository"
	)

func main() {
	driver.DatabaseConnection()
	repository.InitiateEnumTable()
	router.Router()
}


