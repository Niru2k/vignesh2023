package driver

import (
	"fmt"
	"echo/helper"

	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)
var err error

func DatabaseConnection() {
	//connecting to postgres-SQL
	connection:= fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", helper.Host, helper.Port, helper.User, helper.Password,helper.Dbname)
	helper.Db, err = gorm.Open(postgres.Open(connection), &gorm.Config{})
	if err != nil {
		fmt.Println("error in connecting with database")
	}
	fmt.Printf("%s,database connection sucessfull\n",helper.Dbname)
}
	
