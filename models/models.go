package models

type Information struct {
	User_id     uint   `json:"-"                   gorm:"primaryKey;autoIncrement:true"`
	Email       string `json:"email"               gorm:"column:email;type:varchar(35)"`
	Username    string `json:"username"            gorm:"column:user_name;type:varchar(100)"`
	Password    string `json:"password"            gorm:"column:password;type:varchar(100)"`
	Role        string `json:"role"                gorm:"column:role;type:varchar(25)"`
	PhoneNumber string `json:"phone_number"        gorm:"column:phone_number;type:varchar(15)"`
}

type Jobposting struct {
	Job_id      uint   `json:"-"                   gorm:"primaryKey;autoIncrement:true"`
	CompanyName string `json:"company_name"        gorm:"column:company_name;type:varchar(50)"`
	Website     string `json:"website"             gorm:"column:website;type:varchar(50)"`
	JobTitle    string `json:"job_title"           gorm:"column:job_title;type:varchar(50)"`
	JobType     string `json:"job_type"            gorm:"column:job_type;type:varchar(50)"`
	City        string `json:"city"                gorm:"column:city;type:varchar(50)"`
	State       string `json:"state"               gorm:"column:state;type:varchar(50)"`
	Country     string `json:"country"             gorm:"column:country;type:varchar(50)"`
	Email       string `json:"email"               gorm:"column:email;type:varchar(50)"`
	Description string `json:"description"         gorm:"column:description;type:varchar(50)"`
}

type Comments struct {
	Comment_Id uint   `json:"-"                    gorm:"primarykey;autoIncrement:true"`
	Comment    string `json:"comment"              gorm:"column:comment;type:varchar(200)"`
	Interest   bool   `json:"interest"             gorm:"column:interest;type:boolean"`
	User_id    uint   `json:"user_id"              gorm:"type:bigint;references:information(user_id)"`
	Job_id     uint   `json:"job_id"               gorm:"type:bigint;references:jobposting(job_id)"`
}


var SigningKey = []byte("secret")
   