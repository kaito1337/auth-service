package config

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
	"os"
)

type MongoDBConnectionConfig struct {
	Host     string
	Port     int
	Database string
	Username string
	Password string
}

type PostgresConnectionConfig struct {
	Host     string
	Port     int
	Database string
	Username string
	Password string
}

type RabbitConnectionConfig struct {
	Username string
	Password string
	Host     string
	Port     int
}

type Config struct {
	Mongo    MongoDBConnectionConfig  `mapstructure:"mongodb"`
	Postgres PostgresConnectionConfig `mapstructure:"postgres"`
	Rabbit   RabbitConnectionConfig   `mapstructure:"rabbit"`
	Web      WebServerConfig          `mapstructure:"web"`
	App      AppConfig                `mapstructure:"app"`
}

type WebServerConfig struct {
	Port int
}

type AppConfig struct {
	AccessTokenExpirationTimeHours  int    `mapstructure:"access-token-expiration-time-hours"`
	RefreshTokenExpirationTimeHours int    `mapstructure:"refresh-token-expiration-time-hours"`
	TokenSecret                     string `mapstructure:"token-secret"`
}

func InitConfiguration() *Config {
	var C *Config = new(Config)
	loadDefault()
	loadFile()
	viper.Unmarshal(C)
	return C
}

func loadDefault() {

	/*
		MongoDB Connection
	*/

	viper.SetDefault("mongodb.host", os.Getenv("MONGO_HOST"))
	viper.SetDefault("mongodb.port", os.Getenv("MONGO_PORT"))
	viper.SetDefault("mongodb.database", os.Getenv("MONGO_DATABASE"))
	viper.SetDefault("mongodb.username", os.Getenv("MONGO_USER"))
	viper.SetDefault("mongodb.password", os.Getenv("MONGO_PASSWORD"))
	/*
		Postgres Connection
	*/

	viper.SetDefault("postgres.host", os.Getenv("POSTGRES_HOST"))
	viper.SetDefault("postgres.port", os.Getenv("POSTGRES_PORT"))
	viper.SetDefault("postgres.database", os.Getenv("POSTGRES_DATABASE"))
	viper.SetDefault("postgres.username", os.Getenv("POSTGRES_USER"))
	viper.SetDefault("postgres.password", os.Getenv("POSTGRES_PASSWORD"))

	viper.SetDefault("web.port", 8080)

	viper.SetDefault("app.access-token-expiration-time-hours", 24)
	viper.SetDefault("app.refresh-token-expiration-time-hours", 720)
	viper.SetDefault("app.token-secret", "secret")

	/*
		RabbitMQ Connection
	*/

	viper.SetDefault("rabbit.host", os.Getenv("RABBITMQ_HOST"))
	viper.SetDefault("rabbit.port", os.Getenv("RABBITMQ_PORT"))
	viper.SetDefault("rabbit.username", os.Getenv("RABBITMQ_USER"))
	viper.SetDefault("rabbit.password", os.Getenv("RABBITMQ_PASSWORD"))
}

func loadFile() {
	viper.SetConfigName("config")
	viper.AddConfigPath("./config/")
	viper.SetConfigType("yaml")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(fmt.Sprintf("Error reading config file, %s. Use default only.", err))
	}
}
