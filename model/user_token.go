package model

import (
	"errors"
	"github.com/jinzhu/gorm"
)

type UserToken struct {
	ID        uint32 `gorm:"primary_key" json:"id"`
	UserName  string `json:"user_name"`
	UserToken string `json:"user_token"`
}

func (ut UserToken) TableName() string {
	return "tbl_user_token"
}

func (ut UserToken) Update(db *gorm.DB) error {
	userToken := UserToken{UserToken: ut.UserToken, UserName: ut.UserName}
	if db.Model(&userToken).Where("user_name=?", ut.UserName).Updates(&userToken).RowsAffected == 0 {
		return errors.New("not a user")
	}
	return nil
}
