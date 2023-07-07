package repository

import (
	//third party package
	"fmt"

	//user defined package
	"echo/helper"
	"echo/models"
)

func InitiateEnumTable() {
	err := helper.Db.Exec("CREATE TYPE status_enum AS ENUM ('australia', 'india', 'usa','england','dubai','singapore','full time', 'part time', 'freelance','contract')").Error
	if err != nil {
		fmt.Println("error in creating select table")
	}
	helper.Db.AutoMigrate(&models.Select{})
	fmt.Println("select table created successfully")
	selects := []models.Select{
		{Country: "australia", Jobtype: "contract",  JobTitle: "civil engineer"},
		{Country: "australia", Jobtype: "part time", JobTitle: "security"},
		{Country: "australia", Jobtype: "full time", JobTitle: "doctor"},
		{Country: "india",     Jobtype: "contract",  JobTitle: "research scientist"},
		{Country: "india",     Jobtype: "full time", JobTitle: "doctor"},
		{Country: "singapore", Jobtype: "contract",  JobTitle: "constructor"},
		{Country: "singapore", Jobtype: "part time", JobTitle: "security"},
		{Country: "singapore", Jobtype: "full time", JobTitle: "hotel manager"},
		{Country: "usa",       Jobtype: "contract",  JobTitle: "data scientist"},
		{Country: "usa",       Jobtype: "full time", JobTitle: "doctor"},
		{Country: "usa",       Jobtype: "part time", JobTitle: "security"},
	}
	helper.Db.Create(&selects)
	fmt.Println(selects)
	
	//view jobs and job type by mentioning countries
	helper.Db.Where("Country=?", "australia").Find(&selects)
	fmt.Println(selects)
	helper.Db.Where("Country=?", "india").Find(&selects)
	fmt.Println(selects)
	helper.Db.Where("Country=?", "usa").Find(&selects)
    fmt.Println(selects)
}
