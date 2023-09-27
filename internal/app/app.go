package app

import (
	config2 "SergeyProject/config"
	"SergeyProject/internal/repository"
	logger "SergeyProject/pkg/logger"
	"SergeyProject/pkg/postgres"
	"context"
	"github.com/gorilla/mux"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

func Run() {
	logger := logger.GetLogger()
	cfg := config2.ParseConfig(config2.ConfigViper())
	ctx := context.Background()

	connPool := postgres.OpenPoolConnection(ctx, cfg)
	if err := connPool.Ping(ctx); err != nil {
		logger.Error("Unable to ping the database connection", "error", err)
		os.Exit(1)
	}
	postgres.RunMigrationsUp(ctx, cfg)

	storage := repository.NewStorage(connPool)

	router := InitRouter(storage)
	srv := InitHTTPServer(router, cfg)

	StartServer(ctx, srv)
}

func InitRouter(storage *repository.Storage) *mux.Router {

	//router initialization
	router := mux.NewRouter()
	return router
}

func InitHTTPServer(router *mux.Router, cfg *config2.Config) http.Server {
	readTimeoutSecondsCount, _ := strconv.Atoi(cfg.Server.ReadTimeout)
	writeTimeoutSecondsCount, _ := strconv.Atoi(cfg.Server.WriteTimeout)
	idleTimeoutSecondsCount, _ := strconv.Atoi(cfg.Server.IdleTimeout)
	bindAddr := cfg.Server.BindAddr

	srv := http.Server{
		Addr:         bindAddr,
		Handler:      router,
		ReadTimeout:  time.Duration(readTimeoutSecondsCount) * time.Second,
		WriteTimeout: time.Duration(writeTimeoutSecondsCount) * time.Second,
		IdleTimeout:  time.Duration(idleTimeoutSecondsCount) * time.Second,
	}
	return srv
}

func StartServer(ctx context.Context, srv http.Server) {
	logger := logger.GetLogger()

	go func() {
		logger.Info("Starting server...")
		err := srv.ListenAndServe()
		if err != nil {
			logger.Error("Server was stopped", "error", err)
			os.Exit(1)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	signal := <-sigChan
	logger.Info("signal has been recieved", "signal", signal)
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	srv.Shutdown(ctx)
}
