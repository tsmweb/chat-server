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
	hostID     string
	goPoolSize int
	host       string
	port       int
	user       string
	password   string
	dbname     string
	dbschema   string
	serverPort int
	grpcPort   int
	privateKey string
	publicKey  string

	kafkaBootstrapServers string
	kafkaClientID         string
	kafkaServersTopic     string
	kafkaUsersTopic       string
	kafkaNewMessagesTopic string
	kafkaOffMessagesTopic string
	kafkaHostTopic        string
	kafkaErrorsTopic      string
	kafkaGroupID          string
)

func Load(workDir string) {
	err := godotenv.Load(path.Join(workDir, "/.env"))
	if err != nil {
		log.Fatalf("Error loading .env file [%s]", workDir)
	}

	hostID = os.Getenv("HOST_ID")
	goPoolSize, _ = strconv.Atoi(os.Getenv("GOPOOL_SIZE"))
	host = os.Getenv("DB_HOST")
	port, _ = strconv.Atoi(os.Getenv("DB_PORT"))
	user = os.Getenv("DB_USER")
	password = os.Getenv("DB_PASSWORD")
	dbname = os.Getenv("DB_DATABASE")
	dbschema = os.Getenv("DB_SCHEMA")

	serverPort, _ = strconv.Atoi(os.Getenv("SERVER_PORT"))
	grpcPort, _ = strconv.Atoi(os.Getenv("GRPC_PORT"))

	privateKey = workDir + "/config/keys/private-key"
	publicKey = workDir + "/config/keys/public-key.pub"

	kafkaBootstrapServers = os.Getenv("KAFKA_BOOTSTRAP_SERVERS")
	kafkaClientID = os.Getenv("KAFKA_CLIENT_ID")
	kafkaServersTopic = os.Getenv("KAFKA_SERVERS_TOPIC")
	kafkaUsersTopic = os.Getenv("KAFKA_USERS_TOPIC")
	kafkaNewMessagesTopic = os.Getenv("KAFKA_NEW_MESSAGES_TOPIC")
	kafkaOffMessagesTopic = os.Getenv("KAFKA_OFF_MESSAGES_TOPIC")
	kafkaErrorsTopic = os.Getenv("KAFKA_ERRORS_TOPIC")
	kafkaHostTopic = fmt.Sprintf("%s_MESSAGES", hostID)
	kafkaGroupID = fmt.Sprintf("%s_GROUP", hostID)
}

func HostID() string {
	return hostID
}

func GoPoolSize() int {
	return goPoolSize
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

func GRPCPort() int {
	return grpcPort
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

func KafkaBootstrapServers() string {
	return kafkaBootstrapServers
}

func KafkaClientID() string {
	return kafkaClientID
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

func KafkaHostTopic() string {
	return kafkaHostTopic
}

func KafkaGroupID() string {
	return kafkaGroupID
}
