package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Perseverance7/grady/internal/handler"
	"github.com/Perseverance7/grady/internal/models"
	"github.com/Perseverance7/grady/internal/repository"
	"github.com/Perseverance7/grady/internal/service"
	"github.com/Perseverance7/grady/pkg/logger"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func main() {
	logger := logger.NewLogger("app.log")

	if err := godotenv.Load(); err != nil {
		logger.Fatal("error with loading env variables", zap.Error(err))
	}

	db, err := repository.NewPostgresDB(&repository.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   os.Getenv("DB_NAME"),
		SSLMode:  os.Getenv("DB_SSL_MODE"),
	})

	if err != nil {
		logger.Fatal("db connect error", zap.Error(err))
	}

	var secretKey = []byte(os.Getenv("SECRET_KEY"))

	repo := repository.NewRepository(db)
	services := service.NewService(repo, secretKey)
	handlers := handler.NewHandler(logger, services)

	srv := new(models.Server)

	go func() {
		if err := srv.Run(os.Getenv("SERVER_HOST"), handlers.InitRouter()); err != nil && err != http.ErrServerClosed {
			logger.Fatal("error running HTTP server", zap.Error(err))
		}
	}()

	logger.Info("Grady started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	logger.Info("Grady shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("error with shutting down server`", zap.Error(err))
	}

	logger.Info("Shutdown complete")
}
