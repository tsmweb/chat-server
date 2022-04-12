package config

import (
	"log"
	"os"
	"path"
	"strconv"

	"github.com/joho/godotenv"
)

var (
	dbHost                 string
	dbPort                 int
	dbUser                 string
	dbPassword             string
	dbName                 string
	dbSchema               string
	serverPort             int
	keySecureFile          string
	pubSecureFile          string
	certSecureFile         string
	kafkaBootstrapServers  string
	kafkaClientID          string
	kafkaGroupEventTopic   string
	kafkaContactEventTopic string
)

func Load(workDir string) {
	err := godotenv.Load(path.Join(workDir, "/.env"))
	if err != nil {
		log.Fatalf("Error loading .env file [%s]", workDir)
	}

	dbHost = os.Getenv("DB_HOST")
	dbPort, _ = strconv.Atoi(os.Getenv("DB_PORT"))
	dbUser = os.Getenv("DB_USER")
	dbPassword = os.Getenv("DB_PASSWORD")
	dbName = os.Getenv("DB_DATABASE")
	dbSchema = os.Getenv("DB_SCHEMA")

	serverPort, _ = strconv.Atoi(os.Getenv("SERVER_PORT"))

	keySecureFile = workDir + "/config/cert/server.pem"
	pubSecureFile = workDir + "/config/cert/server.pub"
	certSecureFile = workDir + "/config/cert/server.crt"

	kafkaBootstrapServers = os.Getenv("KAFKA_BOOTSTRAP_SERVERS")
	kafkaClientID = os.Getenv("KAFKA_CLIENT_ID")
	kafkaGroupEventTopic = os.Getenv("KAFKA_GROUP_EVENT_TOPIC")
	kafkaContactEventTopic = os.Getenv("KAFKA_CONTACT_EVENT_TOPIC")
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

func KafkaBootstrapServers() string {
	return kafkaBootstrapServers
}

func KafkaClientID() string {
	return kafkaClientID
}

func KafkaGroupEventTopic() string {
	return kafkaGroupEventTopic
}

func KafkaContactEventTopic() string {
	return kafkaContactEventTopic
}
