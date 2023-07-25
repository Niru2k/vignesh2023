package repository

import (
	//user defined package

	"echo/models"

	"gorm.io/gorm"
)

var Db *gorm.DB

func CreateTables() {
	Db.AutoMigrate(&models.Information{})
	Db.AutoMigrate(&models.Jobposting{})
	Db.AutoMigrate(&models.Comments{})
}

func CreateUser(user models.Information) {
	Db.Create(&user)
}

func JobPosting(post models.Jobposting) error {
	err := Db.Create(&post).Error
	return err
}
func ReadUserByEmail(user models.Information) (models.Information, error) {
	err := Db.Where("email=?", user.Email).First(&user).Error
	return user, err
}

func GetAllPosts() ([]models.Jobposting, error) {
	var creates []models.Jobposting
	err := Db.Debug().Find(&creates).Error
	return creates, err
}
func Getjobpostid(companyID string) (models.Jobposting, error) {
	var create models.Jobposting
	err := Db.Where("job_id=?", companyID).First(&create).Error
	return create, err
}

func ReadJobPostById(companyID string) (models.Jobposting, error) {
	var updatedJob models.Jobposting
	err := Db.Where("job_id=?", companyID).First(&updatedJob).Error
	return updatedJob, err
}

func UpdateJob(companyID string, updatedjob models.Jobposting) error {
	err := Db.Where("job_id=?", companyID).Save(&updatedjob).Error
	return err
}

func DeleteJob(companyID string, deletejob models.Jobposting) {
	Db.Where("job_id=?", companyID).Delete(&deletejob)
}

func GetJobpostByCompanyName(companyName string) ([]models.Jobposting, error) {
	var create []models.Jobposting
	err := Db.Where("company_name=?", companyName).Find(&create).Error
	return create, err
}

func CommentPosting(post models.Comments) error {
	err := Db.Create(&post).Error
	return err
}

func GetAllComments() ([]models.Comments, error) {
	var viewcomments []models.Comments
	err := Db.Debug().Find(&viewcomments).Error
	return viewcomments, err
}

func ReadCommentById(commentID string) (models.Comments, error) {
	var getcomment models.Comments
	err := Db.Where("comment_id=?", commentID).First(&getcomment).Error
	return getcomment, err
}

func UpdateComment(commentid string, updatecomment models.Comments) error {
	err := Db.Where("comment_id=?", commentid).Save(&updatecomment).Error
	return err
}

func DeleteComment(CommentID string, deletecomment models.Comments) {
	Db.Where("comment_id=?", CommentID).Delete(&deletecomment)
}
