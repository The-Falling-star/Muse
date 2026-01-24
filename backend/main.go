package main

import (
	"context"
	"errors"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"connectrpc.com/connect"
	"github.com/ling/muse/api"
	"github.com/ling/muse/config"
	"github.com/ling/muse/gen/muse/museconnect"
	"github.com/ling/muse/middleware"

	log "github.com/sirupsen/logrus"
)

func main() {
	// 解析命令行参数
	configPath := flag.String("config", "config.yaml", "配置文件路径")
	flag.Parse()

	// 加载配置
	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	// 初始化数据库连接
	if err := config.InitDatabase(&cfg.Database); err != nil {
		log.Fatalf("初始化数据库失败: %v", err)
	}
	defer func() {
		if err := config.CloseDatabase(); err != nil {
			log.Errorf("关闭数据库连接失败: %v", err)
		}
	}()

	// 创建认证中间件（跳过注册、登录和健康检查接口）
	authInterceptor := middleware.AuthInterceptor([]string{
		museconnect.UserServiceRegisterProcedure, // 注册接口
		museconnect.UserServiceLoginProcedure,    // 登录接口
	})

	// 创建HTTP服务器
	mux := http.NewServeMux()
	mux.Handle(museconnect.NewUserServiceHandler(
		api.NewUserServer(),
		connect.WithInterceptors(authInterceptor),
	))
	mux.Handle(museconnect.NewCharacterServiceHandler(
		api.NewCharacterServer(),
		connect.WithInterceptors(authInterceptor),
	))
	mux.Handle(museconnect.NewChatServiceHandler(
		api.NewChatServer(),
		connect.WithInterceptors(authInterceptor),
	))
	mux.Handle(museconnect.NewPresetServiceHandler(
		api.NewPresetServer(),
		connect.WithInterceptors(authInterceptor),
	))
	mux.Handle(museconnect.NewRegexRuleServiceHandler(
		api.NewRegexRuleServer(),
		connect.WithInterceptors(authInterceptor),
	))
	mux.Handle(museconnect.NewWorldInfoServiceHandler(
		api.NewWorldInfoServer(),
		connect.WithInterceptors(authInterceptor),
	))

	// 健康检查端点
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	server := &http.Server{
		Addr:         cfg.Server.Address(),
		Handler:      mux,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// 在goroutine中启动服务器
	go func() {
		log.Infof("服务器启动成功，监听地址: %s", cfg.Server.Address())
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("服务器启动失败: %v", err)
		}
	}()

	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("正在关闭服务器...")

	// 优雅关闭，等待最多5秒
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Errorf("服务器关闭异常: %v", err)
	}

	log.Info("服务器已关闭")
}
