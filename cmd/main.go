package main

import (
	"filestore-server/global"
	"filestore-server/handler"
	"filestore-server/model"
	"filestore-server/pkg/setting"
	"filestore-server/routers"
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"golang.org/x/net/context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
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

	router := routers.NewRouter()
	s := &http.Server{
		Addr:           ":8080",
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	//http.HandleFunc("/file/upload", handler.UploadHandler)
	//http.HandleFunc("/file/upload/suc", handler.UploadSucHandler)
	//http.HandleFunc("/file/meta", handler.GetFileMetaHandler)
	//http.HandleFunc("/file/download", handler.DownloadHandler)
	//http.HandleFunc("/file/update", handler.FileMetaUpdateHandler)
	//http.HandleFunc("/file/delete", handler.FileDeleteHandler)
	//http.HandleFunc("/user/signup", handler.SignupHandler)
	//http.HandleFunc("/user/signin", handler.SignInHandler)
	//http.HandleFunc("/user/info", handler.UserInfoHandler)
	//todo 要改造注册成为中间件：HTTPInterceptor
	//http.HandleFunc("/user/info", handler.HTTPInterceptor(handler.UserInfoHandler))
	// todo 需要使用router进行gin改造
	//http.HandleFunc("/file/fastupload", handler.HTTPInterceptor(handler.TryFastUploadHandler))
	//todo 分块上传接口-->进行gin改造
	http.HandleFunc("/file/mpupload/init", handler.HTTPInterceptor(handler.InitialMultipartUploadHandler))
	http.HandleFunc("/file/mpupload/uppart", handler.HTTPInterceptor(handler.UploadPartHandler))
	http.HandleFunc("/file/mpupload/complete", handler.HTTPInterceptor(handler.CompleteUploadHandler))
	////这里就不搞配置文件那套了，直接读取端口进行访问
	//err := http.ListenAndServe(":8080", nil)
	//if err != nil {
	//	fmt.Printf("Failed to start server,err:%s", err.Error())
	//}

	go func() {
		err := s.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("s.ListenAndServe err:%v", err)
		}
	}()
	//等待信号中断
	quit := make(chan os.Signal)
	//接受syscall.SiGINT 和syscall.SIGTERM信号
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shuting down server...")

	//最大时间控制，用于通知服务器有5s时间来处理原来的请求
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := s.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown:", err)
	}
	log.Println("Server existing")

}
func init() {
	err := setupSetting()
	if err != nil {
		log.Fatalf("init.setupSetting err: %v", err)
	}
	err = setupFlag()
	if err != nil {
		log.Fatalf("init.setupFlag err: %v", err)
	}
}

// 设置控制参数
func setupSetting() error {
	setting, err := setting.NewSetting(strings.Split(config, ",")...)
	if err != nil {
		return err
	}
	err = setting.ReadSection("Server", &global.ServerSetting)
	if err != nil {
		return err
	}

	err = setting.ReadSection("Database", &global.DatabaseSetting)
	if err != nil {
		return err
	}

	err = setupDBEngine()
	if err != nil {
		log.Fatalf("init.setupDBEngine err: %v", err)
	}

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

// 设置数据库驱动
func setupDBEngine() error {
	var err error
	global.DBEngine, err = model.NewDBEngine(global.DatabaseSetting)
	if err != nil {
		return err
	}

	return nil
}
