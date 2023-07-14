package driver

import (
	"echo/helper"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var err error

func DatabaseConnection() {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println(err)
		return
	}
	host := os.Getenv("HOST")
	port := os.Getenv("PORT")
	password := os.Getenv("PASSWORD")
	dbname := os.Getenv("DBNAME")
	user := os.Getenv("USER")

	//connecting to postgres-SQL
	connection := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	helper.Db, err = gorm.Open(postgres.Open(connection), &gorm.Config{})
	if err != nil {
		fmt.Println("error in connecting with database")
	}
	fmt.Printf("%s,database connection sucessfull\n", helper.Dbname)
}
