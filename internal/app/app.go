package app

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/rogue0026/shortener/internal/config"
	"github.com/rogue0026/shortener/internal/http-server/handlers/urls/add"
	"github.com/rogue0026/shortener/internal/http-server/handlers/urls/get"
	"github.com/rogue0026/shortener/internal/http-server/handlers/urls/ping"
	"github.com/rogue0026/shortener/internal/http-server/handlers/users/login"
	"github.com/rogue0026/shortener/internal/http-server/handlers/users/register"
	"github.com/rogue0026/shortener/internal/http-server/middlewares"
	"github.com/rogue0026/shortener/internal/storage/sqlite"
	"github.com/rogue0026/shortener/pkg/logger"
	"github.com/sirupsen/logrus"
)

func Run() {
	cfg := config.MustLoad()
	log := logger.Init(cfg.Env)
	log.Debug("logger initialized")

	dataStorage, err := sqlite.New(&cfg)
	if err != nil {
		panic(err)
	}

	router := setupRouter(log)

	router.Get("/{short_url}", get.New(log, dataStorage))
	router.Get("/ping", ping.New(dataStorage, log))

	// router.Post("/api/shorten", middlewares.WithCompressing(shorten.New(urlStorage, logger, cfg.Address))) // todo
	router.Post("/", add.New(log, dataStorage, cfg.Address))
	router.Post("/api/users/account", register.New(log, dataStorage))
	router.Post("/api/users/login", login.New(log, dataStorage))

	server := http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		WriteTimeout: cfg.WriteTimeout,
		ReadTimeout:  cfg.ReadTimeout,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		log.Debugf("starting http server listening at %s\n", server.Addr)
		if err := server.ListenAndServe(); err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				log.Debugln("http-server stopped")
			}
		}
	}()

	<-stop

	if err := dataStorage.Close(); err != nil {
		log.Errorln(err.Error())
	} else {
		log.Infoln("data storage successfully closed")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Errorln(err.Error())
	}

	log.Infoln("service stopped")
}

func setupRouter(l *logrus.Logger) *chi.Mux {
	router := chi.NewRouter()
	router.Use(middlewares.WithLogging(l))

	return router
}
