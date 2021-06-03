package setting

import (
	"log"
	"os"
	"path"
	"strconv"

	"github.com/joho/godotenv"
)

var (
	localhost       string
	goPoolSize      int
	host            string
	port            int
	user            string
	password        string
	dbname          string
	dbschema        string
	serverPort      int
	privateKey      string
	publicKey       string
	kafkaBrokerUrls string
)

func Load(workDir string) {
	err := godotenv.Load(path.Join(workDir, "/.env"))
	if err != nil {
		log.Fatalf("Error loading .env file [%s]", workDir)
	}

	localhost = os.Getenv("LOCAL_HOST")
	goPoolSize, _ = strconv.Atoi(os.Getenv("GOPOOL_SIZE"))
	host = os.Getenv("PGHOST")
	port, _ = strconv.Atoi(os.Getenv("PGPORT"))
	user = os.Getenv("PGUSER")
	password = os.Getenv("PGPASSWORD")
	dbname = os.Getenv("PGDATABASE")
	dbschema = os.Getenv("PGSCHEMA")

	serverPort, _ = strconv.Atoi(os.Getenv("SERVER_PORT"))

	privateKey = workDir + "/common/setting/keys/private-key"
	publicKey = workDir + "/common/setting/keys/public-key.pub"

	kafkaBrokerUrls = os.Getenv("KAFKA_BROKER_URLS")
}

func Localhost() string {
	return localhost
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

func KafkaBrokerUrls() string {
	return kafkaBrokerUrls
}
