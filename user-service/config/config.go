package config

import (
	"log"
	"os"
	"path"
	"strconv"

	"github.com/joho/godotenv"
)

var (
	hostID                 string
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
	metricsSendInterval    int
	kafkaBootstrapServers  string
	kafkaClientID          string
	kafkaGroupEventTopic   string
	kafkaContactEventTopic string
	kafkaEventsTopic       string
	kafkaMetricsTopic      string
)

func Load(workDir string) error {
	err := godotenv.Load(path.Join(workDir, "/.env"))
	if err != nil {
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

	metricsSendInterval, err = strconv.Atoi(os.Getenv("METRICS_SEND_INTERVAL")) //sec
	if err != nil {
		return err
	}

	kafkaBootstrapServers = os.Getenv("KAFKA_BOOTSTRAP_SERVERS")
	kafkaClientID = os.Getenv("KAFKA_CLIENT_ID")
	kafkaGroupEventTopic = os.Getenv("KAFKA_GROUP_EVENT_TOPIC")
	kafkaContactEventTopic = os.Getenv("KAFKA_CONTACT_EVENT_TOPIC")
	kafkaEventsTopic = os.Getenv("KAFKA_EVENTS_TOPIC")
	kafkaMetricsTopic = os.Getenv("KAFKA_METRICS_TOPIC")

	return nil
}

func HostID() string {
	return hostID
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

func MetricsSendInterval() int {
	return metricsSendInterval
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

func KafkaEventsTopic() string {
	return kafkaEventsTopic
}

func KafkaMetricsTopic() string {
	return kafkaMetricsTopic
}
