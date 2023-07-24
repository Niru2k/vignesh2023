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
	logs "echo/log"
	"echo/models"
	"echo/repository"

	//third party package
	"github.com/dgrijalva/jwt-go"
	"github.com/fatih/structs"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

// SignUp API
func Signup(c echo.Context) error {
	var user models.Information
	log := logs.Logs()
	log.Info("Signup api called successfully")
	if err := c.Bind(&user); err != nil {
		log.Error("error:'Invalid Format' status:400")
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"Error":  "Invalid Format",
			"status": 400,
		})
	}
	//validates correct email format
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	if !emailRegex.MatchString(user.Email) {
		log.Error("error:'Invalid Email Format' status:400")
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"Error":  "Invalid Email Format",
			"status": 400,
		})
	}
	//make sure username field should not be empty
	if user.Username == "" {
		log.Error("error:'Username field should not be empty' status:400")
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"Error":  "Username field should not be empty",
			"status": 400,
		})
	}
	//password should have minimum 8 character
	if len(user.Password) < 8 {
		log.Error("error:'Password should be more than 8 characters' status:400")
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"Error":  "Password should be more than 8 characters",
			"status": 400,
		})

	}
	//passwords are stored in hashing method in the database
	password, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Error(err)
		return nil
	}
	user.Password = string(password)
	if user.Role == "" {
		log.Error("error:'Role field should not be empty' status:400")
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"Error":  "Role field should not be empty",
			"status": 400,
		})
	}
	// Validate phone number
	phoneNumber := strings.TrimSpace(user.PhoneNumber)
	// Use regular expression to validate numeric characters and length
	phoneRegex := regexp.MustCompile(`^[0-9]{10}$`)
	if !phoneRegex.MatchString(phoneNumber) {
		log.Error("error:'Invalid phone number format' status:400")
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"Error":  "Invalid phone number format",
			"status": 400,
		})
	}
	_, err = repository.ReadUserByEmail(user)
	if err == nil {
		log.Error("error:'user already exist' status:400")
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "user already exist",
			"status":  400,
		})
	}
	repository.CreateUser(user)
	log.Info("message:'sign up successfull' status:200")
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "sign up successfull",
		"status":  200,
	})
}

// Login API
func Login(c echo.Context) error {
	log := logs.Logs()
	log.Info("login api called successfully")
	var login models.Information
	if err := c.Bind(&login); err != nil {
		log.Error("error:'Invalid Format' status:400")
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "Invalid Format",
			"status":  400,
		})
	}
	//verify the email whether its already registered in the SignUp API or not
	verify, err := repository.ReadUserByEmail(login)
	if err == nil {
		//checks whether the given password matches with the email
		if err := bcrypt.CompareHashAndPassword([]byte(verify.Password), []byte(login.Password)); err != nil {
			log.Error("error:'Password Not Matching' status:400")
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"Error":  " Password Not Matching",
				"status": 400,
			})
		}
		//generates token when email and password matches
		login.Role = verify.Role
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"email": login.Email,
			"role":  login.Role,
			"exp":   time.Now().Add(time.Hour * 24).Unix(),
		})
		tokenString, err := token.SignedString(models.SigningKey)
		if err != nil {
			log.Error("error:'Failed To Generate Token' status:400")
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"Error":  "Failed To Generate Token",
				"status": 400,
			})
		}
		log.Info("message:'Login Successful' status:200")
		return c.JSON(http.StatusAccepted, map[string]interface{}{
			"message": "Login Successful",
			"token":   tokenString,
			"status":  200,
		})
	}
	log.Error("error:'login failed' status:400")
	return c.JSON(http.StatusBadRequest, map[string]interface{}{
		"Error":  "login failed",
		"status": 400,
	})
}

// Job Posting API
func Jobposting(c echo.Context) error {
	//allows only admins to post job details
	log := logs.Logs()
	log.Info("job posting api called successfully")
	err := authentication.AdminAuth(c)
	if err != nil {
			log.Error("error:'only admin have the access' status:401")
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"Error":  "only admin have the access",
			"status": 401,
		})
	}
	var post models.Jobposting
	if err := c.Bind(&post); err != nil {
		log.Error("error:'invalid format' status:400")
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"Error":  "invalid format",
			"status": 400,
		})
	}
	//for specifying field name should be empty
	fields := structs.Names(&models.Jobposting{})
	for _, field := range fields {
		if reflect.ValueOf(&post).Elem().FieldByName(field).Interface() == "" {
			check := fmt.Sprintf("missing %s", field)
			log.Error("error:'field details missing' status:400")
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"Error":  check,
				"status": 400,
			})
		}
	}

	//validates correct email format
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	if !emailRegex.MatchString(post.Email) {
		log.Error("error:'invalid Email format' status:400")
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"Error":  "Invalid Email Format.",
			"status": 400,
		})
	}

	err = repository.JobPosting(post)
	if err != nil {
		log.Error("error:'error in creating job posting' status:400")
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"Error":  "error in creating job posting",
			"status": 400,
		})
	}
	log.Info("error:'Job Details Successfully Posted' status:200")
	return c.JSON(http.StatusBadRequest, map[string]interface{}{
		"message": "Job Details Successfully Posted",
		"status":  200,
	})
}

// Get all company job posting details
func GetJobPostingDetails(c echo.Context) error {
	// Allows both admin and user to have access
	log := logs.Logs()
	log.Info("GetJobPostingDetails called successfully")
	err := authentication.CommonAuth(c)
	if err != nil {
		log.Error("error:'only admin and user have the access' status:401")
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"Error":  "only admin and user have the access",
			"status": 401,
		})
	}
	creates, err := repository.GetAllPosts()
	if err != nil {
		log.Error("error:'no record found' status:404")
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"error":  "no record found",
			"status": 404,
		})
	}
	return c.JSON(http.StatusOK, creates)
}

// get jobs posting detail by using job post ID
func GetJobPostingByID(c echo.Context) error {
	log := logs.Logs()
	log.Info("GetJobPostingbyID API  called successfully")
	err := authentication.CommonAuth(c)
	if err != nil {
		log.Error("error:'only admin and user have the access' status:401")
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"Error":  "only admin and user have the access",
			"status": 401,
		})
	}
	companyID := c.Param("id")
	create, err := repository.Getjobpostid(companyID)
	if err != nil {
		log.Error("error:'job post does not exist' status:404")
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"error":  "job post does not exist",
			"status": 404,
		})
	}
	return c.JSON(http.StatusOK, create)
}

// update job posting details by using jobpost ID
func UpdateJob(c echo.Context) error {
	log := logs.Logs()
	log.Info("UpdateJob API called Successfully")
	//allows only admins to update job details
	err := authentication.AdminAuth(c)
	if err != nil {
		log.Error("error:'only admin have the access' status:401")
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"Error":  "only admin have the access",
			"status": 401,
		})
	}
	companyID := c.Param("id")
	updatedjob, err := repository.ReadJobPostById(companyID)
	if err == nil {
		log.Error("error:'can't update status:400")
		if err := c.Bind(&updatedjob); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"error":  "can't update",
				"status": 400,
			})
		}
		//Validates correct email format to be entered
		emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
		if !emailRegex.MatchString(updatedjob.Email) {
			log.Error("error:'Invalid Email Format' status:400")
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"Error":  "Invalid Email Format.",
				"status": 400,
			})
		}
		err := repository.UpdateJob(companyID, updatedjob)
		if err != nil {
			log.Error("error:'job id not found' status:404")
			return c.JSON(http.StatusNotFound, map[string]interface{}{
				"Error":  " job id not found",
				"status": 404,
			})
		}
		log.Info("message:'job updated successfully' status:200")
		return c.JSON(http.StatusOK, map[string]interface{}{
			"message": "job updated successfully",
			"status":  200,
		})
	}
	log.Error("Error:'job post not found' status:404")
	return c.JSON(http.StatusNotFound, map[string]interface{}{
		"Error":  "job post not found",
		"status": 404,
	})
}

// Deletes the jobpost by using jobpost id
func DeleteJob(c echo.Context) error {
	log := logs.Logs()
	log.Info("DeleteJob called successfully")
	//allows only admins to delete job details
	err := authentication.AdminAuth(c)
	if err != nil {
		log.Error("Error:'only admin have the access' status:401")
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"Error":  "only admin have the access",
			"status": 401,
		})
	}
	companyID := c.Param("id")
	deletejob, err := repository.ReadJobPostById(companyID)
	if err == nil {

		repository.DeleteJob(companyID, deletejob)
		log.Info("message:'job post deleted successfully' status:200")
		return c.JSON(http.StatusOK, map[string]interface{}{
			"message": " job post deleted successfully",
			"status":  200,
		})
	}
	log.Error("Error:'job id not found' status:404")
	return c.JSON(http.StatusNotFound, map[string]interface{}{
		"Error":  " job id not found",
		"status": 404,
	})
}

// get all job posted details from a specific company name
func GetJobPostingByCompany(c echo.Context) error {
	log := logs.Logs()
	log.Info("GetJobPostingByCompany API called successfully")
	// Allows both admin and user to have access
	err := authentication.CommonAuth(c)
	if err != nil {
		log.Error("Error:'only admin and user have the access' status:401")
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"Error":  "only admin and user have the access",
			"status": 401,
		})
	}
	company_jobs := c.Param("companyname")
	company_name, err := repository.GetJobpostByCompanyName(company_jobs)
	if err != nil || len(company_name) == 0 {
		log.Error("Error:'company post does not exist' status:404")
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"error":  "company post does not exist",
			"status": 404,
		})
	}
	return c.JSON(http.StatusOK, company_name)
}

// user commenting job post API
func UserComments(c echo.Context) error {
	log := logs.Logs()
	log.Info("UserComments API called successfully")
	//allows only user to post comment
	err := authentication.UserAuth(c)
	if err != nil {
		log.Error("Error:'only user have the access' status:401")
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"Error":  "only user have the access",
			"status": 401,
		})
	}
	var postComments models.Comments
	if err := c.Bind(&postComments); err != nil {
		log.Error("Error:'Invalid Format' status:400")
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"Error":  "Invalid Format",
			"status": 400,
		})
	}
	if postComments.Comment == "" {
		log.Error("Error:'add comment to the post' status:400")
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"Error":  "add comment to the post",
			"status": 400,
		})
	}
	jobId := strconv.Itoa(int(postComments.Job_id))
	_, err = repository.Getjobpostid(jobId)
	if err != nil {
		log.Error("Error:'job ID not found' status:404")
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"Error":  "job ID not found",
			"status": 404,
		})
	}

	err = repository.CommentPosting(postComments)
	if err == nil {
		log.Info("message:'comment posted successfully' status:200")
		return c.JSON(http.StatusOK, map[string]interface{}{
			"message": "comment posted successfully",
			"status":  200,
		})
	}
	log.Error("Error:'Posting a comment failed' status:400")
	return c.JSON(http.StatusBadRequest, map[string]interface{}{
		"Error":  "Posting a comment failed",
		"status": 400,
	})
}

// getting all user comments API
func GetUserComments(c echo.Context) error {
	log := logs.Logs()
	log.Info("GetUserComments API called successfully")
	err := authentication.CommonAuth(c)
	if err != nil {
		log.Error("Error:'only admin and user have the access' status:401")
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"Error":  "only admin and user have the access",
			"status": 401,
		})
	}
	viewcomments, err := repository.GetAllComments()
	if err != nil {
		log.Error("Error:'nothing to see here' status:404")
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"error":  "nothing to see here",
			"status": 404,
		})
	}
	return c.JSON(http.StatusOK, viewcomments)
}

// Get specific comment API
func GetCommentByID(c echo.Context) error {
	log := logs.Logs()
	log.Info("GetCommentByID API called successfully")
	err := authentication.CommonAuth(c)
	if err != nil {
		log.Error("Error:'only admin and user have the access' status:401")
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"Error":  "only admin and user have the access",
			"status": 401,
		})
	}
	var getcomment models.Comments
	commentID := c.Param("id")
	getcomment, err = repository.ReadCommentById(commentID)
	if err != nil {
		log.Error("Error:'no comment found for this id' status:401")
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"error":  "no comment found for this id",
			"status": 404,
		})
	}
	return c.JSON(http.StatusOK, getcomment)
}

// Updating user comment API
func UpdateComment(c echo.Context) error {
	log := logs.Logs()
	log.Info("UpdateComment API called successfully")
	//allows only user to update comment
	err := authentication.UserAuth(c)
	if err != nil {
		log.Error("Error:'only user have the access' status:401")
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"Error":  "only user have the access",
			"status": 401,
		})
	}
	commentid := c.Param("id")
	updatecomment, err := repository.ReadCommentById(commentid)
	if err == nil {
		log.Error("Error:'invalid format' status:400")
		if err := c.Bind(&updatecomment); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"error":  "invalid format",
				"status": 400,
			})
		}
		err := repository.UpdateComment(commentid, updatecomment)
		if err != nil {
			log.Error("Error:'comment id not found ' status:404")
			return c.JSON(http.StatusNotFound, map[string]interface{}{
				"Error":  "comment id not found",
				"status": 404,
			})
		}
		log.Info("message:'comment updated successfully ' status:200")
		return c.JSON(http.StatusOK, map[string]interface{}{
			"message": "comment updated successfully",
			"status":  200,
		})
	}
	log.Error("Error:'comment post not found ' status:404")
	return c.JSON(http.StatusNotFound, map[string]interface{}{
		"Error":  "comment post not found",
		"status": 404,
	})
}

// Deleting user comment API
func DeleteCommentById(c echo.Context) error {
	log := logs.Logs()
	log.Info("DeleteCommentById API called successfully")
	//allows only user to update comment
	err := authentication.UserAuth(c)
	if err != nil {
		log.Error("Error:'only user have the access' status:401")
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"Error":  "only user have the access",
			"status": 401,
		})
	}
	CommentID := c.Param("id")
	deletecomment, err := repository.ReadCommentById(CommentID)
	if err == nil {
		repository.DeleteComment(CommentID, deletecomment)
		log.Info("message:'comment deleted successfully ' status:200")
		return c.JSON(http.StatusOK, map[string]interface{}{
			"message": " comment deleted successfully",
			"status":  200,
		})
	}
	log.Error("Error:'comment id not found' status:404")  
	return c.JSON(http.StatusNotFound, map[string]interface{}{
		"Error":  "comment id not found",
		"status": 404,
	})
}
