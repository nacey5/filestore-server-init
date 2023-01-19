package db

import (
	"github.com/gin-gonic/gin"
)

type UserTokenMeta struct{}

func NewUserTokenMeta() UserTokenMeta {
	return UserTokenMeta{}
}

// UpdateToken 刷新token,这个功能在gorm实现了
func (ut UserTokenMeta) UpdateToken(c *gin.Context) {
}
