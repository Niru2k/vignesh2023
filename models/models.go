package models

type User struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

//using slice to append users while signing up
var Users []User
var SigningKey = []byte("enter secret key")
