package handler

import (
	//built in package
	"fmt"
	"net/http"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

	//user defined package
	"echo/authentication"
	"echo/repository"

	// "echo/middlewares"
	"echo/models"

	//third party package
	"github.com/dgrijalva/jwt-go"
	"github.com/fatih/structs"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"golang.org/x/crypto/bcrypt"
)

// Signing Up API
func Signup(c echo.Context) error {
	var user models.Information
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"warning": "Invalid Format",
			"status":  400,
		})
	}
	//validates correct email format
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	if !emailRegex.MatchString(user.Email) {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"warning": "Invalid Email Format",
			"status":  400,
		})
	}
	if user.Username == "" {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"warning": "Username field should not be empty",
			"status":  400,
		})
	}
	if len(user.Password) < 8 {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"warning": "Password should be more than 8 characters",
			"status":  400,
		})

	}
	password, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Error(err)
		return nil
	}
	user.Password = string(password)
	if user.Role == "" {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"warning": "Role field should not be empty",
			"status":  400,
		})
	}
	// Validate phone number
	phoneNumber := strings.TrimSpace(user.PhoneNumber)
	// Use regular expression to validate numeric characters and length
	phoneRegex := regexp.MustCompile(`^[0-9]{10}$`)
	if !phoneRegex.MatchString(phoneNumber) {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"warning": "Invalid phone number format",
			"status":  400,
		})
	}
	_, err = repository.ReadUserByEmail(user)
	if err == nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "user already exist",
			"status":  400,
		})
	}
	repository.CreateUser(user)
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "sign up successfull",
		"status":  200,
	})
}

// Login API
func Login(c echo.Context) error {
	var login models.Information
	if err := c.Bind(&login); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "Invalid Format",
			"status":  400,
		})
	}
	//verify the email whether its already registered in the SignUp API or not
	verify, err := repository.ReadUserByEmail(login)
	if err == nil {
		//checks whether the given password matches with the email
		if err := bcrypt.CompareHashAndPassword([]byte(verify.Password), []byte(login.Password)); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"warning": " Password Not Matching",
				"status":  400,
			})
		}
		login.Role = verify.Role
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"email": login.Email,
			"role":  login.Role,
			"exp":   time.Now().Add(time.Hour * 24).Unix(),
		})
		tokenString, err := token.SignedString(models.SigningKey)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"warning": "Failed To Generate Token",
				"status":  400,
			})
		}
		return c.JSON(http.StatusAccepted, map[string]interface{}{
			"message": "Login Successful",
			"token":   tokenString,
			"status":  200,
		})
	}
	return c.JSON(http.StatusBadRequest, map[string]interface{}{
		"warning": "login failed",
		"status":  200,
	})
}

// Job Posting API
func Jobposting(c echo.Context) error {
	//allows only admins to post job details
	err := authentication.AdminAuth(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"message": "only admin have the access",
			"status":  401,
		})
	}
	var post models.Jobposting
	if err := c.Bind(&post); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "invalid format",
			"status":  400,
		})
	}

	fields := structs.Names(&models.Jobposting{})
	for _, field := range fields {
		if reflect.ValueOf(&post).Elem().FieldByName(field).Interface() == "" {
			check := fmt.Sprintf("missing %s", field)
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"message": check,
				"status":  400,
			})
		}
	}

	//validates correct email format
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	if !emailRegex.MatchString(post.Email) {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"warning": "Invalid Email Format.",
			"status":  400,
		})
	}

	err = repository.JobPosting(post)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"warning": "error in creating job posting",
			"status":  400,
		})
	}
	return c.JSON(http.StatusBadRequest, map[string]interface{}{
		"SUCCESS": "Job Details Successfully Posted",
		"status":  200,
	})
}

// get all company job posting details
func GetJobPostingDetails(c echo.Context) error {
	err := authentication.CommonAuth(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"message": "only admin and user have the access",
			"status":  401,
		})
	}
	creates, err := repository.GetAllPosts()
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"error":  "nothing to see here",
			"status": 404,
		})
	}
	return c.JSON(http.StatusOK, creates)
}

// get jobs posted by company by using their ID
func GetJobPostingByID(c echo.Context) error {
	err := authentication.CommonAuth(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"message": "only admin and user have the access",
			"status":  401,
		})
	}
	companyID := c.Param("id")
	create, err := repository.Getjobpostid(companyID)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"error":  "job post does not exist",
			"status": 404,
		})
	}
	return c.JSON(http.StatusOK, create)
}

// update job post details by using company id
func UpdateJob(c echo.Context) error {
	//allows only admins to update job details
	err := authentication.AdminAuth(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"message": "only admin have the access",
			"status":  401,
		})
	}
	companyID := c.Param("id")
	updatedjob, err := repository.ReadJobPostById(companyID)
	if err == nil {
		if err := c.Bind(&updatedjob); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"error":  "can't update",
				"status": 400,
			})
		}
		emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
		if !emailRegex.MatchString(updatedjob.Email) {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"warning": "Invalid Email Format.",
				"status":  400,
			})
		}
		err := repository.UpdateJob(companyID, updatedjob)
		if err != nil {
			return c.JSON(http.StatusNotFound, map[string]interface{}{
				"warning": " job id not found",
				"status":  400,
			})
		}
		return c.JSON(http.StatusOK, map[string]interface{}{
			"message": "Job updated successfully",
			"status":  200,
		})
	}
	return c.JSON(http.StatusNotFound, map[string]interface{}{
		"message": "job post not found",
		"status":  404,
	})
}

// DeleteJob handles the DELETE request to delete a job posting by company ID
func DeleteJob(c echo.Context) error {
	//allows only admins to delete job details
	err := authentication.AdminAuth(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"message": "only admin have the access",
			"status":  401,
		})
	}
	companyID := c.Param("id")
	deletejob, err := repository.ReadJobPostById(companyID)
	if err == nil {

		repository.DeleteJob(companyID, deletejob)
		return c.JSON(http.StatusOK, map[string]interface{}{
			"message": " job post deleted successfully",
			"status":  200,
		})
	}
	return c.JSON(http.StatusNotFound, map[string]interface{}{
		"warning": " job id not found",
		"status":  404,
	})
}

// get all job post available
func GetJobPostingByCompany(c echo.Context) error {
	err := authentication.CommonAuth(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"message": "only admin and user have the access",
			"status":  401,
		})
	}
	company_jobs := c.Param("companyname")
	company_name, err := repository.GetJobpostByCompanyName(company_jobs)
	if err != nil || len(company_name) == 0 {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"error":  "company post does not exist",
			"status": 404,
		})
	}
	return c.JSON(http.StatusOK, company_name)
}

// user commenting job post API
func UserComments(c echo.Context) error {
	//allows only user to post comment
	err := authentication.UserAuth(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"message": "only user have the access",
			"status":  401,
		})
	}
	var postComments models.Comments
	if err := c.Bind(&postComments); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"warning": "Invalid Format",
			"status":  400,
		})
	}
	if postComments.Comment == "" {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"warning": "add comment",
			"status":  400,
		})
	}
	jobId := strconv.Itoa(int(postComments.Job_id))
	_, err = repository.Getjobpostid(jobId)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"warning": "job ID not found",
			"status":  404,
		})
	}

	err = repository.CommentPosting(postComments)
	if err == nil {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"message": "comment posted successfully",
			"status":  200,
		})
	}
	return c.JSON(http.StatusBadRequest, map[string]interface{}{
		"message": "Posting a comment failed",
		"status":  400,
	})
}

// getting all user comments API
func GetUserComments(c echo.Context) error {
	err := authentication.CommonAuth(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"message": "only admin and user have the access",
			"status":  401,
		})
	}
	viewcomments, err := repository.GetAllComments()
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"error":  "nothing to see here",
			"status": 404,
		})
	}
	return c.JSON(http.StatusOK, viewcomments)
}

// Get specific comment API
func GetCommentByID(c echo.Context) error {
	err := authentication.CommonAuth(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"message": "only admin and user have the access",
			"status":  401,
		})
	}
	var getcomment models.Comments
	commentID := c.Param("id")
	getcomment, err = repository.ReadCommentById(commentID)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"error":  "no comment found for this id",
			"status": 404,
		})
	}
	return c.JSON(http.StatusOK, getcomment)
}

// Updating user comment API
func UpdateComment(c echo.Context) error {
	//allows only user to update comment
	err := authentication.UserAuth(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"message": "only user have the access",
			"status":  401,
		})
	}
	commentid := c.Param("id")
	updatecomment, err := repository.ReadCommentById(commentid)
	if err == nil {
		if err := c.Bind(&updatecomment); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"error":  "invalid format",
				"status": 400,
			})
		}
		err := repository.UpdateComment(commentid, updatecomment)
		if err != nil {
			return c.JSON(http.StatusNotFound, map[string]interface{}{
				"warning": "comment id not found",
				"status":  404,
			})
		}
		return c.JSON(http.StatusOK, map[string]interface{}{
			"message": "comment updated successfully",
			"status":  200,
		})
	}
	return c.JSON(http.StatusNotFound, map[string]interface{}{
		"message": "comment post not found",
		"status":  404,
	})
}

// Deleting user comment API
func DeleteCommentById(c echo.Context) error {
	//allows only user to update comment
	err := authentication.UserAuth(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"message": "only user have the access",
			"status":  401,
		})
	}
	CommentID := c.Param("id")
	deletecomment, err := repository.ReadCommentById(CommentID)
	if err == nil {
		repository.DeleteComment(CommentID, deletecomment)
		return c.JSON(http.StatusOK, map[string]interface{}{
			"message": " comment deleted successfully",
			"status":  200,
		})
	}
	return c.JSON(http.StatusNotFound, map[string]interface{}{
		"warning": "Invalid comment id",
		"status":  404,
	})
}
