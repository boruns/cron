package main

import (
	"context"
	"crontab/global"
	"crontab/initialize"
	"crontab/services"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"
)

func main() {
	router := initialize.Routers()
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", global.Settings.Port),
		Handler: router,
	}

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	//任务调度协程
	go schedulerJob()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}

	log.Println("Server exiting")
}

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	initConnect()
}

func initConnect() {
	initialize.InitConfig() //初始化配置文件
	initialize.InitLogger()
	initialize.InitMysqlDb()
	initialize.InitRedis()
	initialize.InitMongodb()
	initialize.InitEtcd()
	initialize.InitTrans(global.Settings.Locale)
}

func schedulerJob() {
	//启动任务调度和执行
	services.InitWatchAndScheduleAndExecJob(context.Background())
}
