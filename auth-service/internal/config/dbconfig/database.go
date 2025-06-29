package dbconfig

import (
	"log"
	"os"
  // "reflect"
	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv"
)

type DBConfig struct {
    Host     string
    Port     string
    User     string
    Password string
    Name     string
    SSLMode  string
}

func NewDBConfig() (*DBConfig) {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Can't load .env file")
	}

	var name []string = []string {"DB_HOST", "DB_PORT", "DB_USER", "DB_PASSWORD", "DB_NAME"}
	mas := make([]string, len(name))

	for idx := range name {
		var exists bool
		mas[idx], exists = os.LookupEnv(name[idx])
		if !exists {
			log.Printf("Can't lookup for %s", name[idx])
			return nil
		}
	}

	return &DBConfig {
    Host: mas[0],
    Port: mas[1],
    User: mas[2],
    Password: mas[3],
    Name: mas[4],
    SSLMode: "disable",
  }
}