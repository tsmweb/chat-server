package config

import (
	"fmt"
	"log"
	"os"
	"path"
	"strconv"

	"github.com/joho/godotenv"
)

var (
	hostID                string
	goPoolSize            int
	dbHost                string
	dbPort                int
	dbUser                string
	dbPassword            string
	dbName                string
	dbSchema              string
	redisHost             string
	redisPassword string
	kafkaBootstrapServers string
	kafkaClientID         string
	kafkaGroupID          string
	kafkaServersTopic     string
	kafkaUsersTopic       string
	kafkaNewMessagesTopic string
	kafkaOffMessagesTopic string
	kafkaErrorsTopic      string
	kafkaGroupEventTopic  string
	kafkaHostTopic        string
)

func Load(workDir string) {
	err := godotenv.Load(path.Join(workDir, "/.env"))
	if err != nil {
		log.Fatalf("Error loading .env file [%s]", workDir)
	}

	hostID = os.Getenv("HOST_ID")
	goPoolSize, _ = strconv.Atoi(os.Getenv("GOPOOL_SIZE"))

	dbHost = os.Getenv("DB_HOST")
	dbPort, _ = strconv.Atoi(os.Getenv("DB_PORT"))
	dbUser = os.Getenv("DB_USER")
	dbPassword = os.Getenv("DB_PASSWORD")
	dbName = os.Getenv("DB_DATABASE")
	dbSchema = os.Getenv("DB_SCHEMA")

	redisHost = os.Getenv("REDIS_HOST")
	redisPassword = os.Getenv("REDIS_PASSWORD")

	kafkaBootstrapServers = os.Getenv("KAFKA_BOOTSTRAP_SERVERS")
	kafkaClientID = os.Getenv("KAFKA_CLIENT_ID")
	kafkaGroupID = os.Getenv("KAFKA_GROUP_ID")
	kafkaServersTopic = os.Getenv("KAFKA_SERVERS_TOPIC")
	kafkaUsersTopic = os.Getenv("KAFKA_USERS_TOPIC")
	kafkaNewMessagesTopic = os.Getenv("KAFKA_NEW_MESSAGES_TOPIC")
	kafkaOffMessagesTopic = os.Getenv("KAFKA_OFF_MESSAGES_TOPIC")
	kafkaErrorsTopic = os.Getenv("KAFKA_ERRORS_TOPIC")
	kafkaGroupEventTopic = os.Getenv("KAFKA_GROUP_EVENT_TOPIC")
	kafkaHostTopic = os.Getenv("KAFKA_HOST_TOPIC")
}

func HostID() string {
	return hostID
}

func GoPoolSize() int {
	return goPoolSize
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

func RedisHost() string {
	return redisHost
}

func RedisPassword() string {
	return redisPassword
}

func KafkaBootstrapServers() string {
	return kafkaBootstrapServers
}

func KafkaClientID() string {
	return kafkaClientID
}

func KafkaGroupID() string {
	return kafkaGroupID
}

func KafkaServersTopic() string {
	return kafkaServersTopic
}

func KafkaUsersTopic() string {
	return kafkaUsersTopic
}

func KafkaNewMessagesTopic() string {
	return kafkaNewMessagesTopic
}

func KafkaOffMessagesTopic() string {
	return kafkaOffMessagesTopic
}

func KafkaErrorsTopic() string {
	return kafkaErrorsTopic
}

func KafkaGroupEventTopic() string {
	return kafkaGroupEventTopic
}

func KafkaHostTopic(serverID string) string {
	return fmt.Sprintf(kafkaHostTopic, serverID)
}
