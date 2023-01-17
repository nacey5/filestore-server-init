package v1

import "github.com/gin-gonic/gin"

type FileMeta struct{}

func NewFileMeta() FileMeta {
	return FileMeta{}
}

func (f FileMeta) UploadHandler(c *gin.Context) {

}

func (f FileMeta) UploadSucHandler(c *gin.Context) {

}

func (f FileMeta) GetFileMetaHandler(c *gin.Context) {

}

func (f FileMeta) DownloadHandler(c *gin.Context) {

}

func (f FileMeta) FileMetaUpdateHandler(c *gin.Context) {

}

func (f FileMeta) FileDeleteHandler(c *gin.Context) {

}
