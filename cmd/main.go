package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"

	"git.01.alem.school/bbaktyke/test.project.git"
	hadler "git.01.alem.school/bbaktyke/test.project.git/pkg/handler"
	"git.01.alem.school/bbaktyke/test.project.git/pkg/repository"
	"git.01.alem.school/bbaktyke/test.project.git/pkg/service"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter))
	if err := initConfig(); err != nil {
		logrus.Fatal("error initializing configs:%s", err.Error())
	}
	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("error loading env variable:%s", err.Error())
	}
	db, err := repository.NewPostgresDB(repository.Confiq{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
	})
	if err != nil {
		logrus.Fatalf("failed to initialize db:%s", err.Error())
	}

	repos := repository.NewRepository(db)
	service := service.NewService(repos)
	handlers := hadler.NewHandler(service)

	srv := new(test.Server)
	go func() {
		if err := srv.Run(viper.GetString("port"), handlers.InitRouters()); err != nil {
			logrus.Fatal("error occurred while running http server: ", err)
		}
	}()
	logrus.Println("http://localhost:8080")
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit
	logrus.Println("Problems-Database Shuting Down")
	if err := srv.Shutdown(context.Background()); err != nil {
		logrus.Errorf("Error occured during server shuting down", err.Error())
	}
	if err := db.Close(); err != nil {
		logrus.Errorf("Error occured during db connection closing", err.Error())
	}
}

func initConfig() error {
	viper.AddConfigPath("config")
	viper.SetConfigName("config")
	if err := viper.ReadInConfig(); err != nil {
		return err
	}
	return nil
}
