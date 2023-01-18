package v1

import (
	"encoding/json"
	_const "filestore-server/const"
	"filestore-server/meta"
	"filestore-server/util"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"os"
	"time"
)

type FileMeta struct{}

func NewFileMeta() FileMeta {
	return FileMeta{}
}

func (f FileMeta) UploadHandler(c *gin.Context) {
	if c.Request.Method == _const.GET {
		//返回上传html页面
		data, err := os.ReadFile("./static/view/index.html")
		if err != nil {
			io.WriteString(c.Writer, "internel server error")
			return
		}
		io.WriteString(c.Writer, string(data))
	} else if c.Request.Method == _const.POST {
		c.Request.ParseForm()
		file, head, err := c.Request.FormFile("file")
		if err != nil {
			fmt.Printf("Failed to get data err:%s", err.Error())
			return
		}
		defer file.Close()

		fileMeta := meta.FileMeta{
			FileName: head.Filename,
			Location: "F:\\go_project_all\\filestore-server\\tmp/" + head.Filename,
			UploadAt: time.Now().Format("2006-01-02 15:04:05"),
		}

		//创建文件句柄来接收文件
		newFile, err := os.Create(fileMeta.Location)
		if err != nil {
			fmt.Printf("Failed to create file err:%s", err.Error())
			return
		}
		defer newFile.Close()

		fileMeta.FileSize, err = io.Copy(newFile, file)
		if err != nil {
			fmt.Printf("Failed to save data into file err:%s", err.Error())
			return
		}

		newFile.Seek(0, 0)
		fileMeta.FileSha1 = util.FileSha1(newFile)
		println(fileMeta.FileSha1)
		meta.UpdateFileMeta(fileMeta)

		http.Redirect(c.Writer, c.Request, "/file/upload/suc", http.StatusFound)

	}
}

func (f FileMeta) UploadSucHandler(c *gin.Context) {
	io.WriteString(c.Writer, "Upload finished!!")
}

func (f FileMeta) GetFileMetaHandler(c *gin.Context) {
	r := c.Request
	w := c.Writer
	r.ParseForm()
	filehash := r.Form["filehash"][0]
	fMeta := meta.GetFileMeta(filehash)
	data, err := json.Marshal(fMeta)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(data)
}

func (f FileMeta) DownloadHandler(c *gin.Context) {
	r := c.Request
	w := c.Writer
	r.ParseForm()
	fsha1 := r.Form.Get("filehash")
	fm := meta.GetFileMeta(fsha1)

	file, err := os.Open(fm.Location)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/octect-stream")
	w.Header().Set("Content-Disposition", "attachment;filename=\""+fm.FileName+"\"")
	w.Write(data)
}

func (f FileMeta) FileMetaUpdateHandler(c *gin.Context) {
	r := c.Request
	w := c.Writer
	r.ParseForm()

	opType := r.Form.Get("op")
	fileSha1 := r.Form.Get("filehash")
	newFileName := r.Form.Get("fileName")

	if opType != "0" {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	if r.Method != _const.POST {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	curFileMeta := meta.GetFileMeta(fileSha1)
	curFileMeta.FileName = newFileName
	meta.UpdateFileMeta(curFileMeta)

	data, err := json.Marshal(curFileMeta)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func (f FileMeta) FileDeleteHandler(c *gin.Context) {
	r := c.Request
	w := c.Writer
	r.ParseForm()
	fileSha1 := r.Form.Get("filehash")
	fMeta := meta.GetFileMeta(fileSha1)

	os.Remove(fMeta.Location)
	meta.RemoveFileMeta(fileSha1)

	w.WriteHeader(http.StatusOK)
}
