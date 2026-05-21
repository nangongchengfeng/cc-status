package main

import (
	"log"
	"os"

	"cc-status/server/internal/config"
	"cc-status/server/internal/handler"
	"cc-status/server/internal/repository"
	"cc-status/server/internal/service"
)

func main() {
	cfg, err := config.Load(os.Getenv)
	if err != nil {
		log.Fatalf("load server config: %v", err)
	}

	db, err := repository.OpenDatabase(cfg.SQLitePath)
	if err != nil {
		log.Fatalf("open sqlite database: %v", err)
	}
	if err := repository.InitializeSchema(db); err != nil {
		log.Fatalf("initialize sqlite schema: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("resolve sql database: %v", err)
	}
	defer sqlDB.Close()

	syncHandler := handler.NewSyncHandler(service.NewSyncService(db))
	modelPricingHandler := handler.NewModelPricingHandler(service.NewModelPricingService(db))
	router := handler.NewRouter(cfg.AuthToken, syncHandler.HandleSync, modelPricingHandler)
	if err := router.Run(cfg.ListenAddr); err != nil {
		log.Fatalf("run http server: %v", err)
	}
}
