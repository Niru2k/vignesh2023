package helper

import (
	//third party package
	"gorm.io/gorm"
)

const (
	Host     = "localhost"
	Port     = 5432
	User     = "postgres"
	Password = "password"
	Dbname   = "jobs"
)

var Db *gorm.DB
