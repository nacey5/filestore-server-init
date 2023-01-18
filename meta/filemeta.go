package meta

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

// GetFileMeta 通过fileSha1获得文件
func GetFileMeta(fileSha1 string) FileMeta {
	return fileMetas[fileSha1]
}

// RemoveFileMeta 删除元信息
func RemoveFileMeta(fileSha1 string) {
	delete(fileMetas, fileSha1)
}
