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
	userInfo := v1.NewUserMeta()

	apiv1 := r.Group("/file")
	{
		apiv1.GET("/upload", fileInfo.UploadHandler)
		apiv1.POST("/upload", fileInfo.UploadHandler)
		apiv1.GET("/upload/suc", fileInfo.UploadSucHandler)
		apiv1.GET("/meta", fileInfo.GetFileMetaHandler)
		apiv1.POST("/download", fileInfo.DownloadHandler)
		apiv1.POST("/update", fileInfo.FileMetaUpdateHandler)
		apiv1.DELETE("/delete", fileInfo.FileDeleteHandler)
	}
	apiv1U := r.Group("/user")
	{
		apiv1U.GET("/", userInfo.SignupHandler)
		apiv1U.POST("/", userInfo.SignupHandler)
	}

	return r
}
