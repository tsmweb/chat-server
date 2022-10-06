package config

import (
	"log"
	"os"
	"path"
	"strconv"

	"github.com/joho/godotenv"
)

var (
	hostID                string
	serverPort            int
	dbHost                string
	dbPort                int
	dbUser                string
	dbPassword            string
	dbName                string
	dbSchema              string
	keySecureFile         string
	pubSecureFile         string
	certSecureFile        string
	expireToken           int
	kafkaBootstrapServers string
	kafkaClientID         string
	kafkaEventsTopic      string
)

func Load(workDir string) error {
	if err := godotenv.Load(path.Join(workDir, "/.env")); err != nil {
		log.Fatalf("Error loading .env file [%s]", workDir)
	}

	hostID = os.Getenv("HOST_ID")

	_serverPort, err := strconv.Atoi(os.Getenv("SERVER_PORT"))
	if err != nil {
		return err
	}
	serverPort = _serverPort

	dbHost = os.Getenv("DB_HOST")
	dbPort, err = strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		return err
	}
	dbUser = os.Getenv("DB_USER")
	dbPassword = os.Getenv("DB_PASSWORD")
	dbName = os.Getenv("DB_DATABASE")
	dbSchema = os.Getenv("DB_SCHEMA")

	keySecureFile = workDir + "/config/cert/server.pem"
	pubSecureFile = workDir + "/config/cert/server.pub"
	certSecureFile = workDir + "/config/cert/server.crt"

	expireToken, err = strconv.Atoi(os.Getenv("EXPIRE_TOKEN")) // hour
	if err != nil {
		return err
	}

	kafkaBootstrapServers = os.Getenv("KAFKA_BOOTSTRAP_SERVERS")
	kafkaClientID = os.Getenv("KAFKA_CLIENT_ID")
	kafkaEventsTopic = os.Getenv("KAFKA_EVENTS_TOPIC")

	return nil
}

func HostID() string {
	return hostID
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

func KeySecureFile() string {
	return keySecureFile
}

func PubSecureFile() string {
	return pubSecureFile
}

func CertSecureFile() string {
	return certSecureFile
}

func ExpireToken() int {
	return expireToken
}

func KafkaBootstrapServers() string {
	return kafkaBootstrapServers
}

func KafkaClientID() string {
	return kafkaClientID
}

func KafkaEventsTopic() string {
	return kafkaEventsTopic
}
