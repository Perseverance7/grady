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
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func main() {
	var secretKey = []byte(os.Getenv("SECRET_KEY"))

	logrus.SetFormatter(new(logrus.JSONFormatter))

	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("error with loading env variables: %s", err.Error())
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
		logrus.Fatalf("db connect error %s", err.Error())
	}

	repo := repository.NewRepository(db)
	services := service.NewService(repo, secretKey)
	handlers := handler.NewHandler(services)

	srv := new(models.Server)

	go func() {
		if err := srv.Run(os.Getenv("SERVER_HOST"), handlers.InitRouter()); err != nil && err != http.ErrServerClosed {
			logrus.Fatalf("error running HTTP server: %s", err.Error())
		}
	}()

	logrus.Print("Grady started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	logrus.Print("Grady shutting down")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logrus.Errorf("error occurred on server shutting down: %s", err.Error())
	}

	logrus.Print("Shutdown complete")
}
