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
	"github.com/Perseverance7/grady/pkg/logging"
	"github.com/joho/godotenv"
)

func main() {
	logger := logging.GetLogger()
	logger.Info("start main")

	if err := godotenv.Load(); err != nil {
		logger.Fatalf("error with loading env variables: %s", err.Error())
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
		logger.Fatalf("db connect error %s", err.Error())
	}

	var secretKey = []byte(os.Getenv("SECRET_KEY"))

	repo := repository.NewRepository(db)
	services := service.NewService(repo, secretKey)
	handlers := handler.NewHandler(services)

	srv := new(models.Server)

	go func() {
		if err := srv.Run(os.Getenv("SERVER_HOST"), handlers.InitRouter()); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("error running HTTP server: %s", err.Error())
		}
	}()

	logger.Print("Grady started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	logger.Print("Grady shutting down")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Errorf("error occurred on server shutting down: %s", err.Error())
	}

	logger.Print("Shutdown complete")
}
