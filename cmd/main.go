package main

import (
	"TodoApp"
	"TodoApp/cfg"
	_ "TodoApp/docs"
	"TodoApp/internal/handler"
	"TodoApp/internal/repository"
	"TodoApp/internal/service"
	"context"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// @title Todo App API
// @version 1.0
// description API server for TODO list application

// @host localhost:8080
// @BasePath /

// @securityDefinitions.apiKey ApiKeyAuth
// @in header
// @name Authorization

func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter))
	if err := cfg.InitConfig(); err != nil {
		logrus.Fatalf("error initializating config: %s", err.Error())
	}

	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("error loading .env file: %s", err.Error())
	}

	db, err := repository.NewPostgresDB(repository.Config{
		Host:     viper.GetString("db.host"),
		SSLMode:  viper.GetString("db.sslmode"),
		Password: os.Getenv("DB_PASSWORD"),
		UserName: viper.GetString("db.username"),
		DBName:   viper.GetString("db.name"),
		Port:     viper.GetString("db.port"),
	})

	defer func(db *sqlx.DB) {
		err = db.Close()
		if err != nil {
			logrus.Fatalf("error closing DB: %s", err.Error())
		} else {
			logrus.Infof("DB closed")
		}
	}(db)

	if err != nil {
		logrus.Fatalf("error initializing DB: %s", err.Error())
		return
	}

	repos := repository.NewRepository(db)
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)
	srv := new(TodoApp.Server)

	go func() {
		if err = srv.Run(viper.GetString("port"), handlers.InitRoutes()); err != nil {
			logrus.Fatalf("error occured while running http server: %v", err)
		}
	}()

	logrus.Info("server started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	logrus.Info("server shutting down")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err = srv.Shutdown(ctx); err != nil {
		logrus.Errorf("error shutting down http server: %v", err)
	} else {
		logrus.Info("server stopped")
	}
}
