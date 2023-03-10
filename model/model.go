package model

import (
	"filestore-server/global"
	"filestore-server/pkg/setting"
	"fmt"
	otgorm "github.com/eddycjy/opentracing-gorm"
	"github.com/jinzhu/gorm"
	"time"
)

type Model struct {
	ID       uint32    `gorm:"primary_key" json:"id"`
	CreateAt time.Time `json:"create_at"`
	UpdateAt time.Time `json:"update_at"`
}

type UserModel struct {
	ID         uint32    `gorm:"primary_key" json:"id"`
	SignupAt   time.Time `json:"signup_at"`
	LastActive time.Time `json:"last_active"`
}

func NewDBEngine(databaseSetting *setting.DatabaseSettingS) (*gorm.DB, error) {
	s := "%s:%s@tcp(%s)/%s?charset=%s&parseTime=%t&loc=Local"
	db, err := gorm.Open(databaseSetting.DBType, fmt.Sprintf(s,
		databaseSetting.UserName,
		databaseSetting.Password,
		databaseSetting.Host,
		databaseSetting.DBName,
		databaseSetting.Charset,
		databaseSetting.ParseTime,
	))
	if err != nil {
		return nil, err
	}

	if global.ServerSetting.RunMode == "debug" {
		db.LogMode(true)
	}
	db.SingularTable(true)
	//db.Callback().Create().Replace("gorm:update_time_stamp", updateTimeStampForCreateCallback)
	//db.Callback().Update().Replace("gorm:update_time_stamp", updateTimeStampForUpdateCallback)

	db.DB().SetMaxIdleConns(databaseSetting.MaxIdleConns)
	db.DB().SetMaxOpenConns(databaseSetting.MaxOpenConns)

	otgorm.AddGormCallbacks(db)
	return db, nil
}

// 更新记录的时间hook
func updateTimeStampForUpdateCallback(scope *gorm.Scope) {
	if _, ok := scope.Get("gorm:update_column"); !ok {
		t := time.Now()
		_ = scope.SetColumn("UpdateAt", t.Format("2006-01-02 15:04:05"))
	}
}

// 创建记录的时间hook
func updateTimeStampForCreateCallback(scope *gorm.Scope) {
	if !scope.HasError() {
		nowTime := time.Now().Format("2006:01:02 15:04:05")
		print(nowTime)

		if createTimeField, ok := scope.FieldByName("CreateAt"); ok {
			if createTimeField.IsBlank {
				_ = createTimeField.Set(nowTime)
			}
		}
		if SignupAtField, ok := scope.FieldByName("SignupAt"); ok {
			if SignupAtField.IsBlank {
				_ = SignupAtField.Set(nowTime)
			}
		}
	}
}
