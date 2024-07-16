package server

import (
	"errors"
	"github.com/go-chi/chi/v5"
	orderComposition "github.com/sashaaro/go-musthave-diploma-tpl/internal/composition/order"
	userComposition "github.com/sashaaro/go-musthave-diploma-tpl/internal/composition/user"
	"github.com/sashaaro/go-musthave-diploma-tpl/internal/config"
	sql2 "github.com/sashaaro/go-musthave-diploma-tpl/internal/db/sql"
	"github.com/sashaaro/go-musthave-diploma-tpl/internal/http/middlware/logging"
	jwt2 "github.com/sashaaro/go-musthave-diploma-tpl/pkg/jwt"
	logging2 "github.com/sashaaro/go-musthave-diploma-tpl/pkg/logging"
	"log"
	"net/http"
)

type App struct {
	logger logging2.Logger
	router *chi.Mux
	cfg    *config.Config
}

var (
	ErrNoSQLConnection = errors.New("нет подключения к БД")
)

func New(cfg *config.Config, logger logging2.Logger) (*App, error) {
	jwtClient := jwt2.NewJwt(*cfg.JWTTokenExp, *cfg.JWTSecretKey)
	sql, err := sql2.NewSQL(*cfg.DatabaseURI)
	defer sql.DB.Close()
	if err != nil {
		logger.Error(err)
		return nil, ErrNoSQLConnection
	}
	router := chi.NewRouter()

	logger.Info("Создание userComposite")
	userComposite, err := userComposition.NewUserComposite(cfg, logger, sql.DB, jwtClient)
	if err != nil {
		logger.Fatalf("Ошибка создания userComposite %v", err)
		return nil, err
	}

	logger.Info("Регистрация /user Роутов")
	userComposite.Handler.Register(router)

	logger.Info("Создание orderComposite")
	orderComposite, err := orderComposition.NewOrderComposite(cfg, logger, sql.DB, jwtClient, userComposite.Service)
	if err != nil {
		logger.Fatalf("Ошибка создания orderComposite %v", err)
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
