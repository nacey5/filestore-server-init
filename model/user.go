package model

import (
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
)

type User struct {
	*UserModel
	UserName       string `json:"user_name"`
	UserPwd        string `json:"user_pwd"`
	Email          string `json:"email"`
	Phone          string `json:"phone"`
	EmailValidated string `json:"email_validated"`
	PhoneValidated string `json:"phone_validated"`
	Profile        string `json:"profile"`
	Status         string `json:"status"`
}

func (u User) TableName() string {
	return "tbl_user"
}

// Signup 注册-create
func (u User) Signup(db *gorm.DB) error {
	if count := db.Create(&u).RowsAffected; count != 1 {
		errMsg := fmt.Sprintf("File with hash:%s has benn uploaded before", u.UserName)
		return errors.New(errMsg)
	}
	return nil
}
