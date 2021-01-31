package main

import (
	"context"
	"fmt"
	"kratos/controller"
	"kratos/dao/mysql"
	"kratos/dao/redis"
	"kratos/logger"
	"kratos/pkg/snowflake"
	"kratos/routes"
	"kratos/settings"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/viper"

	"go.uber.org/zap"
)

func main() {
	if err := settings.Init(); err != nil {
		fmt.Printf("init setting failed: %v \n", err)
	}
	if err := logger.Init(); err != nil {
		fmt.Printf("init loggerfailed: %v \n", err)

	}
	defer func() {
		_ = zap.L().Sync()
	}()
	zap.L().Debug("logger init success...")
	if err := mysql.Init(); err != nil {
		fmt.Printf("init mysql failed %v\n", err)
		zap.L().Panic("init mysql failed ", zap.Error(err))
	}
	defer mysql.Close()
	if err := redis.Init(); err != nil {
		fmt.Printf("init redis failed %v\n", err)

		zap.L().Panic("init redis failed ", zap.Error(err))
	}
	defer redis.Close()

	if err := snowflake.Init(viper.GetString("app.start_time"), viper.GetInt64("app.machine_id")); err != nil {
		fmt.Printf("init snowflake failed %v\n", err)
		zap.L().Panic("init snowflake failed ", zap.Error(err))
	}

	if err := controller.InitTrans("zh"); err != nil {
		fmt.Printf("init validator trans failed %v\n", err)
		zap.L().Panic("init validator trans failed ", zap.Error(err))
	}

	r := routes.Setup()
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", viper.GetInt("app.port")),
		Handler: r,
	}

	go func() {
		// 开启一个goroutine启动服务
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			zap.L().Fatal("listen: %s\n", zap.Error(err))
		}
	}()

	// 等待中断信号来优雅地关闭服务器，为关闭服务器操作设置一个5秒的超时
	quit := make(chan os.Signal, 1) // 创建一个接收信号的通道
	// kill 默认会发送 syscall.SIGTERM 信号
	// kill -2 发送 syscall.SIGINT 信号，我们常用的Ctrl+C就是触发系统SIGINT信号
	// kill -9 发送 syscall.SIGKILL 信号，但是不能被捕获，所以不需要添加它
	// signal.Notify把收到的 syscall.SIGINT或syscall.SIGTERM 信号转发给quit
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM) // 此处不会阻塞
	<-quit                                               // 阻塞在此，当接收到上述两种信号时才会往下执行
	zap.L().Info("Shutdown Server ...")
	// 创建一个5秒超时的context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// 5秒内优雅关闭服务（将未处理完的请求处理完再关闭服务），超过5秒就超时退出
	if err := srv.Shutdown(ctx); err != nil {
		zap.L().Fatal("Server Shutdown", zap.Error(err))
	}

	zap.L().Info("Server exiting")
}
