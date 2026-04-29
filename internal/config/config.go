package config

import (
	"fmt"
	"log"
	"tests/internal/models"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Config struct {
	DBHost     string `env:"DB_HOST,required"`
	DBPort     string `env:"DB_PORT,required"`
	DBUser     string `env:"DB_USER,required"`
	DBPassword string `env:"DB_PASSWORD,required"`
	DBName     string `env:"DB_NAME,required"`
	SslMode    string `env:"DB_SSL,required"`
	AppPort    string `env:"APP_PORT,required"`
}

func LoadConfig() *Config {
	_ = godotenv.Load()

	cfg := Config{}

	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("Error parsing config: %v", err)
	}

	return &cfg
}

func (c *Config) GetConnectionString() string {
	return fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		c.DBHost, c.DBPort, c.DBUser, c.DBName, c.DBPassword, c.SslMode)
}

func (c *Config) InitDatabase() *gorm.DB {
	dsn := c.GetConnectionString()
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	err = db.AutoMigrate(&models.Contract{}, &models.Organization{})
	if err != nil {
		panic("failed to migrate database")
	}

	sqlDB, _ := db.DB()
	sqlDB.SetMaxOpenConns(50)
	sqlDB.SetMaxIdleConns(25)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)

	return db
}
