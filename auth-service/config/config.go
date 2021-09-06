package config

import (
	"log"
	"os"
	"path"
	"strconv"

	"github.com/joho/godotenv"
)

var (
	serverPort  int
	dbHost      string
	dbPort      int
	dbUser      string
	dbPassword  string
	dbName      string
	dbSchema    string
	privateKey  string
	publicKey   string
	expireToken int
)

func Load(workDir string) {
	if err := godotenv.Load(path.Join(workDir, "/.env")); err != nil {
		log.Fatalf("Error loading .env file [%s]", workDir)
	}

	dbHost = os.Getenv("DB_HOST")
	dbPort, _ = strconv.Atoi(os.Getenv("DB_PORT"))
	dbUser = os.Getenv("DB_USER")
	dbPassword = os.Getenv("DB_PASSWORD")
	dbName = os.Getenv("DB_DATABASE")
	dbSchema = os.Getenv("DB_SCHEMA")

	serverPort, _ = strconv.Atoi(os.Getenv("SERVER_PORT"))

	privateKey = workDir + "/config/keys/private-key"
	publicKey = workDir + "/config/keys/public-key.pub"
	expireToken, _ = strconv.Atoi(os.Getenv("EXPIRE_TOKEN")) //hour
}

func ServerPort() int {
	return serverPort
}

func DBHost() string {
	return dbHost
}

func DBPort() int {
	return dbPort
}

func DBUser() string {
	return dbUser
}

func DBPassword() string {
	return dbPassword
}

func DBName() string {
	return dbName
}

func DBSchema() string {
	return dbSchema
}

func PathPrivateKey() string {
	return privateKey
}

func PathPublicKey() string {
	return publicKey
}

func ExpireToken() int {
	return expireToken
}
