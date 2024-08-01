package server

import (
	"github.com/go-chi/chi/v5"
	orderComposition "github.com/sashaaro/go-musthave-diploma/internal/composition/order"
	userComposition "github.com/sashaaro/go-musthave-diploma/internal/composition/user"
	"github.com/sashaaro/go-musthave-diploma/internal/config"
	sql2 "github.com/sashaaro/go-musthave-diploma/internal/db/sql"
	"github.com/sashaaro/go-musthave-diploma/internal/http/middlware/logging"
	jwt2 "github.com/sashaaro/go-musthave-diploma/pkg/jwt"
	logging2 "github.com/sashaaro/go-musthave-diploma/pkg/logging"
	"log"
	"net/http"
)

type App struct {
	logger logging2.Logger
	router *chi.Mux
	cfg    *config.Config
}

func New(cfg *config.Config, logger logging2.Logger) (*App, error) {
	jwtClient := jwt2.NewJwt(*cfg.JWTTokenExp, *cfg.JWTSecretKey)
	sql, err := sql2.NewSQL(*cfg.DatabaseURI)
	if err != nil {
		return nil, err
	}
	defer sql.DB.Close()
	router := chi.NewRouter()

	logger.Info("Создание userComposite")
	userComposite, err := userComposition.NewUserComposite(cfg, logger, sql.DB, jwtClient)
	if err != nil {
		return nil, err
	}

	logger.Info("Регистрация /user Роутов")
	userComposite.Handler.Register(router)

	logger.Info("Создание orderComposite")
	orderComposite, err := orderComposition.NewOrderComposite(cfg, logger, sql.DB, jwtClient, userComposite.Service, userComposite.Storage)
	if err != nil {
		return nil, err
	}

	logger.Info("Запуск ProcessingOrders")
	go func() {
		orderComposite.Service.StartProcessingOrders()
	}()

	logger.Info("Регистрация /user/order Роутов")
	orderComposite.Handler.Register(router)

	app := &App{
		logger: logger,
		router: router,
		cfg:    cfg,
	}

	logger.Infof("Start Listen Port %v", *cfg.Port)
	log.Fatal(http.ListenAndServe(*cfg.Port, logging.WithLogging(router, logger)))

	return app, nil
}
