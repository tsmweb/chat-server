package config

import (
	"fmt"
	"log"
	"os"
	"path"
	"runtime"
	"strconv"

	"github.com/joho/godotenv"
)

var (
	hostID                  string
	goPoolSize              int
	dbHost                  string
	dbPort                  int
	dbUser                  string
	dbPassword              string
	dbName                  string
	dbSchema                string
	redisHost               string
	redisPassword           string
	kafkaBootstrapServers   string
	kafkaClientID           string
	kafkaGroupID            string
	kafkaUsersTopic         string
	kafkaUsersPresenceTopic string
	kafkaNewMessagesTopic   string
	kafkaOffMessagesTopic   string
	kafkaGroupEventTopic    string
	kafkaContactEventTopic  string
	kafkaHostTopic          string
	kafkaEventsTopic        string
)

func Load(workDir string) error {
	err := godotenv.Load(path.Join(workDir, "/.env"))
	if err != nil {
		log.Fatalf("Error loading .env file [%s]", workDir)
	}

	hostID = os.Getenv("HOST_ID")
	goPoolSize, err = strconv.Atoi(os.Getenv("GOPOOL_SIZE"))
	if err != nil {
		goPoolSize = runtime.NumCPU()
	}

	dbHost = os.Getenv("DB_HOST")
	dbPort, err = strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		return err
	}
	dbUser = os.Getenv("DB_USER")
	dbPassword = os.Getenv("DB_PASSWORD")
	dbName = os.Getenv("DB_DATABASE")
	dbSchema = os.Getenv("DB_SCHEMA")

	redisHost = os.Getenv("REDIS_HOST")
	redisPassword = os.Getenv("REDIS_PASSWORD")

	kafkaBootstrapServers = os.Getenv("KAFKA_BOOTSTRAP_SERVERS")
	kafkaClientID = os.Getenv("KAFKA_CLIENT_ID")
	kafkaGroupID = os.Getenv("KAFKA_GROUP_ID")
	kafkaUsersTopic = os.Getenv("KAFKA_USERS_TOPIC")
	kafkaUsersPresenceTopic = os.Getenv("KAFKA_USERS_PRESENCE_TOPIC")
	kafkaNewMessagesTopic = os.Getenv("KAFKA_NEW_MESSAGES_TOPIC")
	kafkaOffMessagesTopic = os.Getenv("KAFKA_OFF_MESSAGES_TOPIC")
	kafkaGroupEventTopic = os.Getenv("KAFKA_GROUP_EVENT_TOPIC")
	kafkaContactEventTopic = os.Getenv("KAFKA_CONTACT_EVENT_TOPIC")
	kafkaHostTopic = os.Getenv("KAFKA_HOST_TOPIC")
	kafkaEventsTopic = os.Getenv("KAFKA_EVENTS_TOPIC")

	return nil
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

func KafkaUsersTopic() string {
	return kafkaUsersTopic
}

func KafkaUsersPresenceTopic() string {
	return kafkaUsersPresenceTopic
}

func KafkaNewMessagesTopic() string {
	return kafkaNewMessagesTopic
}

func KafkaOffMessagesTopic() string {
	return kafkaOffMessagesTopic
}

func KafkaGroupEventTopic() string {
	return kafkaGroupEventTopic
}

func KafkaContactEventTopic() string {
	return kafkaContactEventTopic
}

func KafkaHostTopic(serverID string) string {
	return fmt.Sprintf("%s_%s", serverID, kafkaHostTopic)
}

func KafkaEventsTopic() string {
	return kafkaEventsTopic
}
