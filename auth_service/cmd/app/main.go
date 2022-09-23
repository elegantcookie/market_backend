package main

import (
	_ "auth_service/docs"
	"auth_service/internal/app"
	"auth_service/internal/config"
	"auth_service/pkg/logging"
	"log"
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
