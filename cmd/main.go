package main

import (
	"filestore-server/global"
	"filestore-server/handler"
	"filestore-server/pkg/setting"
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"time"
)

var (
	port         string
	runMode      string
	config       string
	isVersion    bool
	buildTime    string
	buildVersion string
	gitCommitID  string
)

func main() {

	if isVersion {
		fmt.Printf("build_time: %s\n", buildTime)
		fmt.Printf("build_version: %s\n", buildVersion)
		fmt.Printf("git_commit_id: %s\n", gitCommitID)
		return
	}
	gin.SetMode(global.ServerSetting.RunMode)

	http.HandleFunc("/file/upload", handler.UploadHandler)
	http.HandleFunc("/file/upload/suc", handler.UploadSucHandler)
	http.HandleFunc("/file/meta", handler.GetFileMetaHandler)
	http.HandleFunc("/file/download", handler.DownloadHandler)
	http.HandleFunc("/file/update", handler.FileMetaUpdateHandler)
	http.HandleFunc("/file/delete", handler.FileDeleteHandler)
	//这里就不搞配置文件那套了，直接读取端口进行访问
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Printf("Failed to start server,err:%s", err.Error())
	}
}

// 设置控制参数
func setupSetting() error {
	setting, err := setting.NewSetting(strings.Split(config, ",")...)
	if err != nil {
		return err
	}
	err = setting.ReadSection("Server", &global.ServerSetting)

	global.ServerSetting.ReadTimeout *= time.Second
	global.ServerSetting.WriteTimeout *= time.Second

	if port != "" {
		global.ServerSetting.HttpPort = port
	}
	if runMode != "" {
		global.ServerSetting.RunMode = runMode
	}
	return nil
}

// 设置标记参数
func setupFlag() error {
	flag.StringVar(&port, "port", "", "启动端口")
	flag.StringVar(&runMode, "runMode", "", "启动模式")
	flag.StringVar(&config, "config", "configs/", "指定要使用的配置文件路径")
	flag.BoolVar(&isVersion, "version", false, "编译信息")
	flag.Parse()
	return nil
}
