package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"

	"git.01.alem.school/bbaktyke/test.project.git"
	"git.01.alem.school/bbaktyke/test.project.git/cache"
	hadler "git.01.alem.school/bbaktyke/test.project.git/internal/handler"
	"git.01.alem.school/bbaktyke/test.project.git/internal/repository"
	"git.01.alem.school/bbaktyke/test.project.git/internal/service"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

func main() {
	// initialize configuration
	logrus.SetFormatter(new(logrus.JSONFormatter))
	if err := initConfig(); err != nil {
		logrus.Fatal("error initializing configs:%s", err.Error())
	}
	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("error loading env variable:%s", err.Error())
	}
	// initialize database
	db, err := connectToDatabase()
	if err != nil {
		logrus.Fatalf("failed to initialize db:%s", err.Error())
	}
	// create RabbitMQ connection
	conn, err := amqp.Dial("amqp://guest:guest@rabbitmq:5672/")
	if err != nil {
		log.Fatalf("failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	// create RabbitMQ channel
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("failed to open a channel: %v", err)
	}
	defer ch.Close()
	// declare a RabbitMQ queue
	q, err := ch.QueueDeclare(
		"my-queue", // name
		false,      // durable
		false,      // delete when unused
		false,      // exclusive
		false,      // no-wait
		nil,        // arguments
	)
	if err != nil {
		log.Fatalf("failed to declare a queue: %v", err)
	}

	repos := repository.NewRepository(db)
	// service := service.NewService(repos)
	authService := service.NewAuthService(repos.Authorization)
	problemService := service.NewProblemService(repos.Problem)
	cache := cache.NewRedisCache("localhost:6379", 0, 10)
	handlers := hadler.NewHandler(authService, problemService, cache, ch, q.Name)
	router := mux.NewRouter()
	handlers.InitAuthRouters(router)
	handlers.InitProblemRouters(router)

	// create server
	srv := new(test.Server)
	go func() {
		if err := srv.Run(viper.GetString("port"), router); err != nil {
			log.Fatalf("error occurred while running http server: %v", err)
		}
	}()
	log.Println("http://localhost:8080")
	log.Printf("server started on port %v", viper.GetString("port"))
	// handle signals
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	sig := <-sigCh
	logrus.Printf("received signal %v, shutting down server", sig)

	// gracefully shutdown server
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logrus.Printf("error occurred while shutting down server: %v", err)
	}

	// close RabbitMQ channel
	if err := ch.Close(); err != nil {
		log.Printf("error occurred while closing RabbitMQ channel: %v", err)
	}

	// close RabbitMQ connection
	if err := conn.Close(); err != nil {
		log.Printf("error occurred while closing RabbitMQ connection: %v", err)
	}
	// close database
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

func connectToDatabase() (*sqlx.DB, error) {
	return repository.NewPostgresDB(repository.Confiq{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
	})
}
