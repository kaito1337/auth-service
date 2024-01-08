package main

import (
	"auth-backend/internal/config"
	"auth-backend/internal/mongodb"
	"auth-backend/internal/postgres"
	"auth-backend/internal/rabbitmq"
	"auth-backend/internal/services"
	"auth-backend/internal/web"
	"auth-backend/internal/web/controllers"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
	"os/signal"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal(fmt.Sprintf("Error loading .env file: %s", err))
	}
	cfg := config.InitConfiguration()
	if cfg == nil {
		log.Fatal("failed to load config")
	}
	mongoClient := mongodb.NewMongoDB(&cfg.Mongo)

	mongoError := mongoClient.Connect()
	if mongoError != nil {
		log.Fatal("failed to connect mongodb:", mongoError)
	}
	pgClient := postgres.NewPostgresDB(&cfg.Postgres)
	pgError := pgClient.Connect()
	if pgError != nil {
		log.Fatal("failed to connect postgres:", pgError)
	}
	migrateError := pgClient.Migrate()
	if migrateError != nil {
		log.Fatal("failed to migrate postgres:", migrateError)
	}

	rmq := rabbitmq.NewRabbitMQ(&cfg.Rabbit)
	err := rmq.Connect()
	if err != nil {
		log.Fatal("failed to connect rabbitmq:", err)
	}

	wApp := web.NewWebServer(&cfg.Web)
	wApp.LogMiddleware(rmq.Channel)
	registerRoutes(&cfg.App, wApp, pgClient)

	go wApp.Run()
	go rmq.Consume()

	defer func() {
		mongoClient.Release()
		rmq.Close()
	}()
	//
	// Gracefully shutdown
	//
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt, os.Kill)
	takeSig := <-sigChan
	log.Print("Shutting down gracefully, ", "signal:", takeSig.String())

}

func registerRoutes(cfg *config.AppConfig, wApp *web.WebServer, pgClient *postgres.PostgresDB) {

	userService := services.NewUserService(pgClient.GetDB())
	authService := services.NewAuthService(cfg, userService)

	wApp.RegisterRoutes([]controllers.Controller{
		controllers.NewAuthController(authService),
		controllers.NewRegisterController(userService),
	})
}
