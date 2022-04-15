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
	serverPort     int
	dbHost         string
	dbPort         int
	dbUser         string
	dbPassword     string
	dbName         string
	dbSchema       string
	maxUploadSize  int64
	keySecureFile  string
	pubSecureFile  string
	certSecureFile string
	filePath       string
	userFilePath   string
	groupFilePath  string
	mediaFilePath  string
)

func Load(workDir string) error {
	if err := godotenv.Load(path.Join(workDir, "/.env")); err != nil {
		log.Fatalf("Error loading .env file [%s]", workDir)
	}

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

	filePath = filepath.Join(workDir, "files")
	userFilePath = filepath.Join(filePath, "user")
	groupFilePath = filepath.Join(filePath, "group")
	mediaFilePath = filepath.Join(filePath, "media")

	if err := os.MkdirAll(userFilePath, os.ModePerm); err != nil {
		return err
	}

	if err := os.MkdirAll(groupFilePath, os.ModePerm); err != nil {
		return err
	}

	if err := os.MkdirAll(mediaFilePath, os.ModePerm); err != nil {
		return err
	}

	dbHost = os.Getenv("DB_HOST")
	dbPort, _ = strconv.Atoi(os.Getenv("DB_PORT"))
	dbUser = os.Getenv("DB_USER")
	dbPassword = os.Getenv("DB_PASSWORD")
	dbName = os.Getenv("DB_DATABASE")
	dbSchema = os.Getenv("DB_SCHEMA")

	keySecureFile = workDir + "/config/cert/server.pem"
	pubSecureFile = workDir + "/config/cert/server.pub"
	certSecureFile = workDir + "/config/cert/server.crt"

	return nil
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

func FilePath() string {
	return filePath
}

func UserFilePath() string {
	return userFilePath
}

func GroupFilePath() string {
	return groupFilePath
}

func MediaFilePath() string {
	return mediaFilePath
}
