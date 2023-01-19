package model

import (
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
)

type File struct {
	ID       uint32 `gorm:"primary_key" json:"id"`
	FileSha1 string `json:"file_sha1"`
	FileName string `json:"file_name"`
	FileSize int64  `json:"file_size"`
	FileAddr string `json:"file_addr"`
	status   int64  `json:"status"`
}

func (f File) TableName() string {
	return "tbl_file"
}

func (f File) Update(db *gorm.DB) error {
	//更新
	if count := db.Create(&f).RowsAffected; count != 1 {
		errMsg := fmt.Sprintf("File with hash:%s has benn uploaded before", f.FileSha1)
		return errors.New(errMsg)
	}
	return nil
}

// Get 通过gorm查询数据库获得文件元数据
func (f File) Get(db *gorm.DB) (File, error) {
	var file File
	err := db.Where("file_sha1= ? and status=1", f.FileSha1).First(&file).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return file, err
	}

	return file, nil
}
