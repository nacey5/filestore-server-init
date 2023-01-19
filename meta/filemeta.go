package meta

import (
	"filestore-server/global"
	"filestore-server/model"
)

// FileMeta 文件元信息结构
type FileMeta struct {
	FileSha1 string
	FileName string
	FileSize int64
	Location string
	UploadAt string
}

// key-->指定为fileSha1为唯一标识
var fileMetas map[string]FileMeta

func init() {
	println("fileMetas初始化")
	fileMetas = make(map[string]FileMeta)
}

// UpdateFileMeta 新增或者更新文件元信息
func UpdateFileMeta(fmeta FileMeta) {
	println("更新map信息")
	fileMetas[fmeta.FileSha1] = fmeta
}

// 更新/新增文件元信息到mysql中
func UpdateFileMetaDB(fmeta FileMeta) bool {
	file := model.File{
		FileSha1: fmeta.FileSha1,
		FileSize: fmeta.FileSize,
		FileName: fmeta.FileName,
		FileAddr: fmeta.Location,
	}
	err := file.Update(global.DBEngine)
	if err == nil {
		return true
	}
	return false
}

// GetFileMeta 通过fileSha1获得文件
func GetFileMeta(fileSha1 string) FileMeta {
	return fileMetas[fileSha1]
}

// GetFileMetaDB 从mysql获取文件元信息
func GetFileMetaDB(fileSha1 string) (FileMeta, error) {
	file := model.File{FileSha1: fileSha1}
	//tfile, err := mydb.GetFileMeta(fileSha1)
	file, err := file.Get(global.DBEngine)
	if err != nil {
		return FileMeta{}, err
	}
	fmeta := FileMeta{
		FileName: file.FileName,
		FileSha1: file.FileSha1,
		FileSize: file.FileSize,
		Location: file.FileAddr}

	return fmeta, nil
}

// RemoveFileMeta 删除元信息
func RemoveFileMeta(fileSha1 string) {
	delete(fileMetas, fileSha1)
}
