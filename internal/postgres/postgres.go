package postgres

import (
	"auth-backend/internal/config"
	"auth-backend/internal/models"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type PostgresDB struct {
	client *gorm.DB
	cfg    *config.PostgresConnectionConfig
}

func NewPostgresDB(cfg *config.PostgresConnectionConfig) *PostgresDB {
	return &PostgresDB{
		client: nil,
		cfg:    cfg,
	}
}

func (p *PostgresDB) Connect() error {
	DSN := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", p.cfg.Host, p.cfg.Port, p.cfg.Username, p.cfg.Password, p.cfg.Database)
	if DBData, err := gorm.Open(postgres.Open(DSN), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	}); err != nil {
		return err
	} else {
		p.client = DBData
	}
	fmt.Println("You successfully connected to PostgreSQL!")
	return nil
}

func (p *PostgresDB) Migrate() error {
	err := p.client.AutoMigrate(&models.User{})
	if err != nil {
		return err
	}
	fmt.Println("Migrated your deployment!")
	return nil
}

func (p *PostgresDB) GetDB() *gorm.DB {
	return p.client
}
