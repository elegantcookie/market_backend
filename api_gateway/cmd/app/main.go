package main

import (
	"api_gateway/internal/app"
	"api_gateway/internal/config"
	"api_gateway/pkg/logging"
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
