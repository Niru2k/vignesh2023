package driver

import (
	//user defined package

	"echo/helper"
	logs "echo/log"
	"echo/repository"

	//built in package
	"fmt"
	"os"

	//third party package

	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var err error

func DatabaseConnection() error {
	log := logs.Logs()
	err := helper.Configure(".env")
	if err != nil {
		log.Error("error in loading env file")
		fmt.Println("error is loading env file ")
	}
	host := os.Getenv("HOST")
	port := os.Getenv("PORT")
	password := os.Getenv("PASSWORD")
	dbname := os.Getenv("DBNAME")
	user := os.Getenv("USER")

	//connecting to postgres-SQL
	connection := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	repository.Db, err = gorm.Open(postgres.Open(connection), &gorm.Config{})
	if err != nil {
		log.Error("error in connecting with database")
		fmt.Println("error in connecting with database")
	}
	log.Info("database connection sucessfull")
	fmt.Printf("%s,database connection sucessfull\n", dbname)
	return nil
}
