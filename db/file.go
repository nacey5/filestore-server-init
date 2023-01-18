package db

import (
	"database/sql"
	mydb "filestore-server/db/mysql"
	"fmt"
)

// OnFileUploadFinished 文件上传完成，保存本地数据库
func OnFileUploadFinished(filehash, filename, fileaddr string, filesize int64) bool {
	conn := mydb.DBConn()
	stmt, err := conn.Prepare(
		"insert ignore into tbl_file (`file_sha1`,`file_name`,`file_size`,`file_addr`,`status`) values (?,?,?,?,1)")
	defer stmt.Close()
	if err != nil {
		fmt.Println("Failed to prepare statement,err:%s", err)
		return false
	}
	ret, err := stmt.Exec(filehash, filename, filesize, fileaddr)
	if err != nil {
		fmt.Println("Failed to Exec sql,err:%s", err)
		return false
	}

	if rf, err := ret.RowsAffected(); nil == err {
		if rf <= 0 {
			fmt.Printf("File with hash:%s has benn uploaded before", filehash)
		}
		return true
	}
	return false
}

type TableFile struct {
	FileHash string
	FileName sql.NullString
	FileSize int64
	FileAddr sql.NullString
}

// GetFileMeta 从mysql获取源文件信息
func GetFileMeta(filehash string) (*TableFile, error) {
	prepare, err := mydb.DBConn().Prepare(
		"select file_sha1,file_addr,file_name,file_size from tbl_file where file_sha1=? and status=1 limit 1")
	defer prepare.Close()
	if err != nil {
		fmt.Println("Failed to prepare statement,err:%s", err)
		return nil, err
	}
	tfile := TableFile{}
	err = prepare.QueryRow(filehash).Scan(&tfile.FileHash, &tfile.FileAddr, &tfile.FileName, &tfile.FileSize)
	if err != nil {
		fmt.Printf(err.Error())
		return nil, err
	}

	return &tfile, nil

}
