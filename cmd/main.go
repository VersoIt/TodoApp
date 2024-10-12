package main

import (
	"TodoApp"
	"TodoApp/cfg"
	"TodoApp/internal/handler"
	"TodoApp/internal/repository"
	"TodoApp/internal/service"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
)

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
	if err = srv.Run(viper.GetString("port"), handlers.InitRoutes()); err != nil {
		logrus.Fatalf("error occured while running http server: %v", err)
	}
}
