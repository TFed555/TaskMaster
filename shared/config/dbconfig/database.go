package dbconfig

import (
	"fmt"
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
)

type DBConfig struct {
    Host     string
    Port     string
    User     string
    Password string
    Name     string
    SSLMode  string
}

func NewDBConfig() (DBConfig) {
	if err := godotenv.Load("../.env", ".env.local"); err != nil {
		log.Fatal("Can't load .env file")
	}

	var name []string = []string {"DB_HOST", "DB_PORT", "DB_USER", "DB_PASSWORD", "DB_NAME"}
	mas := make([]string, len(name))

	for idx := range name {
		var exists bool
		mas[idx], exists = os.LookupEnv(name[idx])
		if !exists {
			log.Printf("Can't lookup for %s", name[idx])
			return DBConfig{}
		}
	}

	return DBConfig {
    Host: mas[0],
    Port: mas[1],
    User: mas[2],
    Password: mas[3],
    Name: mas[4],
    SSLMode: "disable",
  }
}

func (dbconf *DBConfig) CreateConnection() (string){
	conn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
			dbconf.User,
			dbconf.Password,
			dbconf.Host,
			dbconf.Port,
			dbconf.Name,
			dbconf.SSLMode)
	return conn
}

func (dbconf *DBConfig) CreateDSN() (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		dbconf.Host,
		dbconf.Port,
		dbconf.User,
		dbconf.Name,
		dbconf.Password,
		dbconf.SSLMode))
	return db, err
}