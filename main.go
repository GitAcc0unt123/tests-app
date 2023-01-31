package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"tests_app/internal/config"
	"tests_app/internal/handler"
	"tests_app/internal/repository"
	"tests_app/internal/server"
	"tests_app/internal/service"
	"tests_app/pkg/hash"
	"tests_app/pkg/token"

	_ "github.com/lib/pq"
)

// @title       API Server
// @version     0.1.0
// @description API Server

// @host     localhost:8080
// @BasePath /api
// @accept   json
// @produce  json
// @schemes  http

// @securityDefinitions.apikey ApiKeyAuth
// @in                         header
// @name                       Authorization

func main() {
	config, err := config.Init("./")
	if err != nil {
		log.Fatalf("fatal error config file: %s", err.Error())
	}

	// wait mongo container
	time.Sleep(3 * time.Second)

	// connect to postgres
	db, err := repository.NewPostgresDB(repository.Config{
		Host:     config.DB.Host,
		Port:     config.DB.Port,
		Username: config.DB.Username,
		DBName:   config.DB.DBName,
		SSLMode:  config.DB.SSLMode,
		Password: config.DB.Password,
	})
	if err != nil {
		log.Fatalf("failed to initialize postgres: %s", err.Error())
	}

	hasher := hash.NewSHA512Hasher(config.PasswordSalt)
	tokenManager := token.NewManager(config.JWT.AccessTokenTTL, config.JWT.SigningKey)

	repos := repository.New(db)
	services := service.New(repos, tokenManager, hasher)
	handlers := handler.New(services, *config)

	srv := new(server.Server)
	c := handlers.InitApiRoutes()
	handlers.InitFileRoutes(c)
	go func() {
		if err := srv.Run(config.Server, c); err != nil {
			log.Fatalf("error occured while running http server: %s", err.Error())
		}
	}()

	log.Print("App Started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	log.Print("App Shutting Down")

	if err := srv.Shutdown(context.Background()); err != nil {
		log.Printf("error occured on server shutting down: %s", err.Error())
	}

	if err := db.Close(); err != nil {
		log.Printf("error occured on postgres connection close: %s", err.Error())
	}
}
