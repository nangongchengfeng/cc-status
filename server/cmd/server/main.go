package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"cc-status/server/internal/config"
	"cc-status/server/internal/handler"
	"cc-status/server/internal/repository"
	"cc-status/server/internal/service"
)

func main() {
	log.SetFlags(log.LstdFlags)
	log.Println("[启动] CC Status Server 启动中...")

	cfg, err := config.Load(os.Getenv)
	if err != nil {
		log.Fatalf("[错误] 配置加载失败: %v", err)
	}
	log.Printf("[配置] 加载完成, 监听地址: %s, 数据库路径: %s", cfg.ListenAddr, cfg.SQLitePath)

	db, err := repository.OpenDatabase(cfg.SQLitePath)
	if err != nil {
		log.Fatalf("[错误] 数据库连接失败: %v", err)
	}
	log.Println("[数据库] 连接成功")

	if err := repository.InitializeSchema(db); err != nil {
		log.Fatalf("[错误] 数据库表结构初始化失败: %v", err)
	}
	log.Println("[数据库] 表结构初始化完成")

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("[错误] 获取数据库实例失败: %v", err)
	}
	defer sqlDB.Close()

	syncHandler := handler.NewSyncHandler(service.NewSyncService(db))
	modelPricingHandler := handler.NewModelPricingHandler(service.NewModelPricingService(db))
	statsHandler := handler.NewStatsHandler(service.NewStatsService(db))
	logsHandler := handler.NewLogsHandler(service.NewLogsService(db))
	log.Println("[服务] 模块初始化完成")

	router := handler.NewRouter(cfg.AuthToken, syncHandler.HandleSync, modelPricingHandler, statsHandler, logsHandler)

	srv := &http.Server{
		Addr:    cfg.ListenAddr,
		Handler: router,
	}

	go func() {
		log.Printf("[启动] 服务已启动, 监听地址: %s", cfg.ListenAddr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("[错误] 服务启动失败: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	for sig := range quit {
		switch sig {
		case syscall.SIGHUP:
			log.Println("[重载] 收到 SIGHUP 信号, 触发配置重载")
		case syscall.SIGINT, syscall.SIGTERM:
			log.Printf("[关闭] 收到 %s 信号, 正在优雅关闭服务...", sig)
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			if err := srv.Shutdown(ctx); err != nil {
				log.Printf("[关闭] 服务关闭超时: %v", err)
			}
			log.Println("[关闭] 服务已安全关闭")
			return
		}
	}
}