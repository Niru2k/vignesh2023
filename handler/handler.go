package handler

import (
	//built in package
	"net/http"
	"regexp"
	"strings"
	"time"

	//user defined package
	"echo/authentication"
	"echo/helper"

	// "echo/middlewares"
	"echo/models"

	//third party package
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
)

// Signing Up API
func Signup(c echo.Context) error {
	var user models.Information
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"warning": "Invalid Format"})
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
	err := helper.Db.Where("email=?", user.Email).First(&user).Error
	if err == nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "user already exist",
			"status":  400,
		})
	}
	helper.Db.Create(&user)
	return c.JSON(http.StatusBadRequest, map[string]interface{}{
		"message": "sign up successfull",
		"status":  200,
	})
}

// Login API
func Login(c echo.Context) error {
	var login models.Information
	var verify models.Information
	if err := c.Bind(&login); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "Invalid Format",
			"status":  400,
		})
	}
	//verify the email whether its already registered in the SignUp API or not
	err := helper.Db.Where("email=?", login.Email).First(&verify).Error
	login.Role = verify.Role
	if err == nil {
		//checks whether the given password matches with the email
		if verify.Password != login.Password {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"warning": " Password Not Matching",
				"status":  400,
			})
		}
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
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "only admin have the access",
			"status":  400,
		})
	}
	var post models.Jobposting
	if err := c.Bind(&post); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "invalid format",
			"status":  400,
		})
	}
	if post.CompanyName == "" || post.Website == "" || post.JobTitle == "" || post.JobType == "" || post.City == "" || post.State == "" || post.Country == "" || post.Email == "" || post.Description == "" {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"warning": "Field Or Field Details Missing",
			"status":  400,
		})
	}
	//validates correct email format
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	if !emailRegex.MatchString(post.Email) {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"warning": "Invalid Email Format.",
			"status":  400,
		})
	}
	helper.Db.Create(&post)
	return c.JSON(http.StatusBadRequest, map[string]interface{}{
		"SUCCESS": "Job Details Successfully Posted",
		"status":  400,
	})
}

// get all company job posting details
func GetJobPostingDetails(c echo.Context) error {
	var creates []models.Jobposting
	err := helper.Db.Debug().Find(&creates).Error
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error":  "nothing to see here",
			"status": 400,
		})
	}
	return c.JSON(http.StatusOK, creates)
}

// get jobs posted by company by using their ID
func GetJobPostingByID(c echo.Context) error {
	var create models.Jobposting
	companyID := c.Param("id")
	err := helper.Db.Where("job_id=?", companyID).First(&create).Error
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error":  "job post does not exist",
			"status": 400,
		})
	}
	return c.JSON(http.StatusOK, create)
}

// update job post details by using company id
func UpdateJob(c echo.Context) error {
	//allows only admins to update job details
	err := authentication.AdminAuth(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "only admin have the access",
			"status":  400,
		})
	}
	companyID := c.Param("id")
	var updatedJob models.Jobposting
	err = helper.Db.Where("job_id=?", companyID).First(&updatedJob).Error
	if err == nil {
		if err := c.Bind(&updatedJob); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"error":  "can't update",
				"status": 400,
			})
		}
		emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
		if !emailRegex.MatchString(updatedJob.Email) {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"warning": "Invalid Email Format.",
				"status":  400,
			})
		}
		err = helper.Db.Where("job_id=?", companyID).Save(&updatedJob).Error
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"warning": "Invalid job id",
				"status":  400,
			})
		}
		return c.JSON(http.StatusOK, map[string]interface{}{
			"message": "Job updated successfully",
			"status":  200,
		})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "job post not found",
		"status":  400,
	})
}

// DeleteJob handles the DELETE request to delete a job posting by company ID
func DeleteJob(c echo.Context) error {
	//allows only admins to delete job details
	err := authentication.AdminAuth(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "only admin have the access",
			"status":  400,
		})
	}
	companyID := c.Param("id")
	var deletejob models.Jobposting
	err = helper.Db.Where("job_id=?", companyID).First(&deletejob).Error
	if err == nil {
		err = helper.Db.Where("job_id=?", companyID).Delete(&deletejob).Error
		return c.JSON(http.StatusOK, map[string]interface{}{
			"message": " job post deleted successfully",
			"status":  200,
		})
	}
	return c.JSON(http.StatusBadRequest, map[string]interface{}{
		"warning": "Invalid job id",
		"status":  400,
	})
}

// get all job post available
func GetJobPostingByCompany(c echo.Context) error {
	var company_name []models.Jobposting
	company_jobs := c.Param("companyname")
	err := helper.Db.Where("company_name=?", company_jobs).Find(&company_name).Error
	if err != nil || len(company_name) == 0 {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error":  "company post does not exist",
			"status": 400,
		})
	}
	return c.JSON(http.StatusOK, company_name)
}

// user commenting job post API
func UserComments(c echo.Context) error {
	//allows only user to post comment
	err := authentication.UserAuth(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "only user have the access",
			"status":  400,
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
	var job models.Jobposting
	result := helper.Db.First(&job, postComments.Job_id)
	if result.Error != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"warning": "job ID not found",
			"status":  400,
		})
	}
	helper.Db.Create(&postComments)
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "comment posted successfully",
		"status":  200,
	})
}

// getting all user comments API
func GetUserComments(c echo.Context) error {
	var viewcomments []models.Comments
	err := helper.Db.Debug().Find(&viewcomments).Error
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error":  "nothing to see here",
			"status": 400,
		})
	}
	return c.JSON(http.StatusOK, viewcomments)
}

// Get specific comment API
func GetCommentByID(c echo.Context) error {
	var getcomment models.Comments
	commentID := c.Param("id")
	err := helper.Db.Where("comment_id=?", commentID).First(&getcomment).Error
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error" :  "no comment found for this id",
			"status": 400,
		})
	}
	return c.JSON(http.StatusOK, getcomment)
}

// Updating user comment API
func UpdateComment(c echo.Context) error {
	//allows only user to update comment
	err := authentication.UserAuth(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "only user have the access",
			"status" :  400,
		})
	}
	commentid := c.Param("id")
	var updatecomment models.Comments
	err = helper.Db.Where("comment_id=?", commentid).First(&updatecomment).Error
	if err == nil {
		if err := c.Bind(&updatecomment); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"error" :  "invalid format",
				"status": 400,
			})
		}
		err = helper.Db.Where("comment_id=?", commentid).Save(&updatecomment).Error
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"warning": "Invalid comment id",
				"status" :  400,
			})
		}
		return c.JSON(http.StatusOK, map[string]interface{}{
			"message": "comment updated successfully",
			"status":  400,
		})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "comment post not found",
		"status":  200,
	})
}

// Deleting user comment API
func DeleteCommentById(c echo.Context) error {
	//allows only user to update comment
	err := authentication.UserAuth(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "only user have the access",
			"status":  400,
		})
	}
	CommentID := c.Param("id")
	var deletecomment models.Comments
	err = helper.Db.Where("comment_id=?", CommentID).First(&deletecomment).Error
	if err == nil {
		err = helper.Db.Where("comment_id=?", CommentID).Delete(&deletecomment).Error
		return c.JSON(http.StatusOK, map[string]interface{}{
			"message": " comment deleted successfully",
			"status":  200,
		})
	}
	return c.JSON(http.StatusBadRequest, map[string]interface{}{
		"warning": "Invalid comment id",
		"status":  400,
	})
}
