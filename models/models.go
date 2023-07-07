package models

type User struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type Jobposting struct {
	CompanyID   string `json:"id"`
	CompanyName string `json:"company name"`
	Website     string `json:"website"`
	JobTitle    string `json:"job title"`
	JobType     string `json:"job type"`
	City        string `json:"city"`
	State       string `json:"state"`
	Country     string `json:"country"`
	Email       string `json:"email"`
	Description string `json:"description"`
}

// implementing enum method
type StatusEnum string
type Select struct {
	Country  StatusEnum `gorm:"type:status_enum"`
	Jobtype  StatusEnum `gorm:"type:status_enum"`
	JobTitle string
}

// using slice method to append users while signing up
var Users []User

//using slice to add job posting
var Job []Jobposting
var SigningKey = []byte("enter secret key")
