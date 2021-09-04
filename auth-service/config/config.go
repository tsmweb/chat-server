package config

import (
	"log"
	"os"
	"path"
	"strconv"

	"github.com/joho/godotenv"
)

var (
	host         string
	port         int
	user         string
	password     string
	dbname       string
	dbschema     string
	serverPort   int
	privateKey   string
	publicKey    string
	expireToken  int
)

func Load(workDir string) {
	if err := godotenv.Load(path.Join(workDir, "/.env")); err != nil {
		log.Fatalf("Error loading .env file [%s]", workDir)
	}

	host = os.Getenv("DB_HOST")
	port, _ = strconv.Atoi(os.Getenv("DB_PORT"))
	user = os.Getenv("DB_USER")
	password = os.Getenv("DB_PASSWORD")
	dbname = os.Getenv("DB_DATABASE")
	dbschema = os.Getenv("DB_SCHEMA")

	serverPort, _ = strconv.Atoi(os.Getenv("SERVER_PORT"))

	privateKey = workDir + "/config/keys/private-key"
	publicKey = workDir + "/config/keys/public-key.pub"

	expireToken = 24 //hour
}

func PathPrivateKey() string {
	return privateKey
}

func PathPublicKey() string {
	return publicKey
}

func ServerPort() int {
	return serverPort
}

func DBHost() string {
	return host
}

func DBPort() int {
	return port
}

func DBUser() string {
	return user
}

func DBPassword() string {
	return password
}

func DBName() string {
	return dbname
}

func DBSchema() string {
	return dbschema
}

func ExpireToken() int {
	return expireToken
}
