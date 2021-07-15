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
	kafkaClientsTopic     string
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
	host = os.Getenv("PGHOST")
	port, _ = strconv.Atoi(os.Getenv("PGPORT"))
	user = os.Getenv("PGUSER")
	password = os.Getenv("PGPASSWORD")
	dbname = os.Getenv("PGDATABASE")
	dbschema = os.Getenv("PGSCHEMA")

	serverPort, _ = strconv.Atoi(os.Getenv("SERVER_PORT"))
	grpcPort, _ = strconv.Atoi(os.Getenv("GRPC_PORT"))

	privateKey = workDir + "/config/keys/private-key"
	publicKey = workDir + "/config/keys/public-key.pub"

	kafkaBootstrapServers = os.Getenv("KAFKA_BOOTSTRAP_SERVERS")
	kafkaClientID = os.Getenv("KAFKA_CLIENT_ID")
	kafkaServersTopic = os.Getenv("KAFKA_SERVERS_TOPIC")
	kafkaClientsTopic = os.Getenv("KAFKA_CLIENTS_TOPIC")
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

func KafkaClientsTopic() string {
	return kafkaClientsTopic
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
