package config

import (
	"os"
  "log"
	_"github.com/joho/godotenv"
  "fmt"
)

type DBConfig struct {
    Host     string
    Port     string
    User     string
    Password string
    Name     string
    SSLMode  string
}

func getEnv(value string) string {
  p, err := os.LookupEnv(value)
  if !err{
      log.Fatal("Unable to load env var")
      return ""
    }
  return p
}

func NewDBConfig() (*DBConfig, error) {
  // if err := godotenv.Load(); err != nil {
	// 	log.Fatal("Can't load .env file")
	// }

  cfg := &DBConfig{
        Host:     getEnv("DB_HOST"),
        Port:     getEnv("DB_PORT"),
        User:     getEnv("DB_USER"),
        Password: getEnv("DB_PASSWORD"),
        Name:     getEnv("DB_NAME"),
        SSLMode:  getEnv("DB_SSLMODE"),
    }
    
    return cfg, nil
}

func (d *DBConfig) Connection() string {
    return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
    d.User, d.Password, d.Host, d.Port, d.Name, d.SSLMode)
}
