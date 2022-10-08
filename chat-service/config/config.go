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
	serverPort              int
	keySecureFile           string
	pubSecureFile           string
	certSecureFile          string
	kafkaBootstrapServers   string
	kafkaClientID           string
	kafkaServersTopic       string
	kafkaUsersTopic         string
	kafkaUsersPresenceTopic string
	kafkaNewMessagesTopic   string
	kafkaOffMessagesTopic   string
	kafkaHostTopic          string
	kafkaGroupID            string
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
	serverPort, err = strconv.Atoi(os.Getenv("SERVER_PORT"))
	if err != nil {
		return err
	}

	keySecureFile = workDir + "/config/cert/server.pem"
	pubSecureFile = workDir + "/config/cert/server.pub"
	certSecureFile = workDir + "/config/cert/server.crt"

	kafkaBootstrapServers = os.Getenv("KAFKA_BOOTSTRAP_SERVERS")
	kafkaClientID = os.Getenv("KAFKA_CLIENT_ID")
	kafkaServersTopic = os.Getenv("KAFKA_SERVERS_TOPIC")
	kafkaUsersTopic = os.Getenv("KAFKA_USERS_TOPIC")
	kafkaUsersPresenceTopic = os.Getenv("KAFKA_USERS_PRESENCE_TOPIC")
	kafkaNewMessagesTopic = os.Getenv("KAFKA_NEW_MESSAGES_TOPIC")
	kafkaOffMessagesTopic = os.Getenv("KAFKA_OFF_MESSAGES_TOPIC")
	kafkaHostTopic = fmt.Sprintf("%s_MESSAGES", hostID)
	kafkaGroupID = os.Getenv("KAFKA_GROUP_ID")
	kafkaEventsTopic = os.Getenv("KAFKA_EVENTS_TOPIC")

	return nil
}

func HostID() string {
	return hostID
}

func GoPoolSize() int {
	return goPoolSize
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

func KafkaUsersPresenceTopic() string {
	return kafkaUsersPresenceTopic
}

func KafkaNewMessagesTopic() string {
	return kafkaNewMessagesTopic
}

func KafkaOffMessagesTopic() string {
	return kafkaOffMessagesTopic
}

func KafkaHostTopic() string {
	return kafkaHostTopic
}

func KafkaGroupID() string {
	return kafkaGroupID
}

func KafkaEventsTopic() string {
	return kafkaEventsTopic
}
