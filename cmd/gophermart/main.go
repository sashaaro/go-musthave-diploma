package main

import (
	server "github.com/sashaaro/go-musthave-diploma/internal"
	"github.com/sashaaro/go-musthave-diploma/internal/config"
	"github.com/sashaaro/go-musthave-diploma/pkg/logging"
	"log"
)

func main() {
	cfg := config.NewConfig()
	err := cfg.Load()
	if err != nil {
		log.Fatalln(err)
		return
	}

	logging.Init()
	logger := logging.NewMyLogger()

	logger.Info("Запуск приложения")
	_, err = server.New(cfg, logger)
	if err != nil {
		logger.Fatalf("Запуск приложения провалено: %s", err)
	}
}
