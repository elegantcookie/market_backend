package main

import (
	"log"
	_ "user_service/docs"
	"user_service/internal/app"
	"user_service/internal/config"
	"user_service/pkg/logging"
)

func main() {
	log.Print("config initialization")
	cfg := config.GetConfig()

	log.Printf("logging initialized.")

	logger := logging.GetLogger(cfg.AppConfig.LogLevel)

	a, err := app.NewApp(cfg, &logger)
	if err != nil {
		logger.Fatal(err)
	}

	logger.Println("running Application")
	a.Run()
}
