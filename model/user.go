package model

import (
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
)

type User struct {
	ID             uint32 `gorm:"primary_key" json:"id"`
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
	u.EmailValidated = "0"
	u.PhoneValidated = "0"
	u.Status = "1"
	if count := db.Create(&u).RowsAffected; count != 1 {
		errMsg := fmt.Sprintf("File with hash:%s has benn uploaded before", u.UserName)
		return errors.New(errMsg)
	}
	return nil
}

// Signin 登陆校验
func (u User) Signin(db *gorm.DB) error {
	var user User
	err := db.Where("user_name=? and user_pwd=?", u.UserName, u.UserPwd).First(&user).Error
	if err != nil {
		return err
	}
	return nil
}
