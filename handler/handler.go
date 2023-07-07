package handler

import (
	//built in package

	"net/http"
	"regexp"

	//user defined package
	"echo/models"

	//third party package
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
)

// for user signing up
func Signup(c echo.Context) error {
	var user models.User
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"warning": "Invalid Format"})
	}


	//make sure that email, username, and password should not be given as an empty field
	if user.Email == ""  {
		return c.JSON(http.StatusBadRequest, map[string]string{"warning": "Email field should not be empty"})
	}
	if user.Username == ""{
		return c.JSON(http.StatusBadRequest, map[string]string{"warning": "Username field should not be empty"})
	}
	if len(user.Password) < 8 {
		return c.JSON(http.StatusBadRequest, map[string]string{"warning": "Password should be more than 8 characters"})
	}
	//validates correct email format
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	if !emailRegex.MatchString(user.Email) {
		return c.JSON(http.StatusBadRequest, map[string]string{"warning": "Invalid Email Format."})
	}

	// Check if the Email already exist are not in Signup  API
	for _, u := range models.Users {
		if u.Email == user.Email {
			return c.JSON(http.StatusBadRequest, map[string]string{"warning": "Email Already Exists"})
		}
	}

	// Generate JWT token while signing up
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email":    user.Email,
		"username": user.Username,
	})
	tokenString, err := token.SignedString(models.SigningKey)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"warning": "Failed To Generate Token"})
	}

	// Saves n number of users to the list
	models.Users = append(models.Users, user)
	return c.JSON(http.StatusOK, map[string]string{"token": tokenString})
}

// user Login API after Signing Up
func Login(c echo.Context) error {
	var user models.User
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid Format"})
	}

	//verify the email whether its already registered in Sign Up API or not
	for _, u := range models.Users {
		if u.Email == user.Email {

			//verify whether the username exist or not in Sign Up API
			if u.Username != user.Username {
				return c.JSON(http.StatusBadRequest, map[string]string{"warning": "Username Does Not Exists "})
			}

			//checks whether the given password matches with the username
			if u.Password != user.Password {
				return c.JSON(http.StatusBadRequest, map[string]string{"warning": " Password Not Matching"})
			}
			return c.JSON(http.StatusAccepted, map[string]string{"warning": "Login Successful"})
		}
	}
	return c.JSON(http.StatusBadRequest, map[string]string{"warning": "login failed"})
}

// API For Job Posting Details
func Jobposting(c echo.Context) error {
	var post models.Jobposting
	if err := c.Bind(&post); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Field Is Misssing"})
	}
	for _, u := range models.Job {
		if u.CompanyID == post.CompanyID {
			return c.JSON(http.StatusBadRequest, map[string]string{"warning": "User Id Already Exist, Please Try Different Id"})
		}
		if u.CompanyName == post.CompanyName {
			return c.JSON(http.StatusBadRequest, map[string]string{"warning": "Company Name Already Exist"})
		}
		if u.Website == post.Website {
			return c.JSON(http.StatusBadRequest, map[string]string{"warning": "Website Already Registered"})
		}
		if u.Email == post.Email {
			return c.JSON(http.StatusBadRequest, map[string]string{"warning": "Email Already Exists"})
		}
	}

	//the field details should not be empty
	if post.CompanyID == "" || post.CompanyName == "" || post.Website == "" || post.JobTitle == "" || post.JobType == "" || post.City == "" || post.State == "" || post.Country == "" || post.Email == "" || post.Description == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"warning": "Field Or Field Details Missing"})
	}

	//validates correct email format
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	if !emailRegex.MatchString(post.Email) {
		return c.JSON(http.StatusBadRequest, map[string]string{"warning": "Invalid Email Format."})
	}

	// Perform  validation and save the new job posting
	models.Job = append(models.Job, post)
	return c.JSON(http.StatusBadRequest, map[string]string{"SUCCESS": "Job Details Successfully Posted"})
}

// get all company job posting details
func GetJobPostingDetails(c echo.Context) error {
	return c.JSON(http.StatusOK, models.Job)
}

// get available jobs posted by company by using their ID
func GetJobPostingByID(c echo.Context) error {
	companyID := c.Param("id")

	// Find all job postings with the matching company ID
	var foundJobPostings []models.Jobposting
	for _, job := range models.Job {
		if job.CompanyID == companyID {
			foundJobPostings = append(foundJobPostings, job)
		}
	}

	// Check if any job postings were found
	if len(foundJobPostings) == 0 {
		return c.JSON(http.StatusNotFound, map[string]string{"WARNING": " Company ID Does Not Exists "})
	}
	return c.JSON(http.StatusOK, foundJobPostings)
}

// update job posting details by using id
func UpdateJob(c echo.Context) error {
	// Parse the company ID from the URL parameter
	companyID := c.Param("id")
	var updatedJob models.Jobposting
	if err := c.Bind(&updatedJob); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	if updatedJob.CompanyID == "" || updatedJob.CompanyName == "" || updatedJob.Website == "" || updatedJob.JobTitle == "" || updatedJob.JobType == "" || updatedJob.City == "" || updatedJob.State == "" || updatedJob.Country == "" || updatedJob.Email == "" || updatedJob.Description == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"warning": "Fields should not be empty"})
	}
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	if !emailRegex.MatchString(updatedJob.Email) {
		return c.JSON(http.StatusBadRequest, map[string]string{"warning": "Invalid Email Format."})
	}
	index := -1
	for i, job := range models.Job {
		if job.CompanyID == companyID {
			index = i
			break
		}
	}

	// If the job with the specified company ID is found, update it
	if index != -1 {
		models.Job[index] = updatedJob
	} else {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Job not found"})
	}

	// Return a success response
	return c.JSON(http.StatusOK, map[string]string{"message": "Job updated successfully"})
}

// DeleteJob handles the DELETE request to delete a job posting by company ID
func DeleteJob(c echo.Context) error {
	// Parse the company ID from the URL parameter
	companyID := c.Param("id")

	// Find the index of the job with the specified company ID in the Job slice
	index := -1
	for i, job := range models.Job {
		if job.CompanyID == companyID {
			index = i
			break
		}
	}

	// If the job with the specified company ID is found, delete it
	if index != -1 {
		// Remove the job from the slice by creating a new slice without the job
		models.Job = append(models.Job[:index], models.Job[index+1:]...)
	} else {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Job not found"})
	}

	// Return a success response
	return c.JSON(http.StatusOK, map[string]string{"message": "Job deleted successfully"})
}
