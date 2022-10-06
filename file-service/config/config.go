package config

import (
	"log"
	"os"
	"path"
	"path/filepath"
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
	maxUploadSize         int64
	keySecureFile         string
	pubSecureFile         string
	certSecureFile        string
	fileDir               string
	userFileDir           string
	groupFileDir          string
	mediaFileDir          string
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

	_maxUploadSize, err := strconv.Atoi(os.Getenv("MAX_UPLOAD_SIZE"))
	if err != nil {
		return err
	}
	maxUploadSize = int64(1024 * 1024 * _maxUploadSize)

	fileDir = filepath.Join(workDir, "files")
	userFileDir = filepath.Join(fileDir, "user")
	groupFileDir = filepath.Join(fileDir, "group")
	mediaFileDir = filepath.Join(fileDir, "media")

	if err = os.MkdirAll(userFileDir, os.ModePerm); err != nil {
		return err
	}

	if err = os.MkdirAll(groupFileDir, os.ModePerm); err != nil {
		return err
	}

	if err = os.MkdirAll(mediaFileDir, os.ModePerm); err != nil {
		return err
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

	keySecureFile = workDir + "/config/cert/server.pem"
	pubSecureFile = workDir + "/config/cert/server.pub"
	certSecureFile = workDir + "/config/cert/server.crt"

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

func MaxUploadSize() int64 {
	return maxUploadSize
}

func SetMaxUploadSize(size int64) {
	maxUploadSize = 1024 * size
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

func FileDir() string {
	return fileDir
}

func UserFileDir() string {
	return userFileDir
}

func GroupFileDir() string {
	return groupFileDir
}

func MediaFileDir() string {
	return mediaFileDir
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
