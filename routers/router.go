package routers

import (
	"filestore-server/global"
	v1 "filestore-server/routers/api/v1"
	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	r := gin.New()
	if global.ServerSetting.RunMode == "debug" {
		r.Use(gin.Logger())
		r.Use(gin.Recovery())
	}

	fileInfo := v1.NewFileMeta()

	apiv1 := r.Group("/api/v1")
	{
		apiv1.GET("/file/upload", fileInfo.UploadHandler)
		apiv1.POST("/file/upload/suc", fileInfo.UploadSucHandler)
		apiv1.POST("/file/meta", fileInfo.GetFileMetaHandler)
		apiv1.POST("/file/download", fileInfo.DownloadHandler)
		apiv1.POST("/file/update", fileInfo.FileMetaUpdateHandler)
		apiv1.DELETE("/file/delete", fileInfo.FileDeleteHandler)
	}

	return r
}
