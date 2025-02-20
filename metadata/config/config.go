package config

import "os"

type Config struct {
	Domain        string
	Port          string
	ConsulAddress string
	KafkaAddress  string
	DbDriver      string
	DbConnection  string
	/* DbHost        string
	DbPort        string
	DbUser        string
	DbPassword    string
	DbName        string
	DbSslMode     string */
}

func NewConfig() *Config {
	config := Config{}
	/* flag.StringVar(&config.Domain, "domain", "localhost", "")
	flag.StringVar(&config.Port, "port", ":8080", "Listen for requests on this port")
	flag.StringVar(&config.ConsulAddress, "consul_address", "localhost:8500", "Register for service discovery")
	flag.StringVar(&config.KafkaAddress, "kafka_address", "localhost:9092", "Publish and consume messages")
	flag.StringVar(&config.DbDriver, "db_driver", "postgres", "Driver for database connection")
	flag.StringVar(&config.DbHost, "db_host", "host.docker.internal", "Database host")
	flag.StringVar(&config.DbPort, "db_port", "5432", "Database port")
	flag.StringVar(&config.DbUser, "db_user", "admin", "Database user")
	flag.StringVar(&config.DbPassword, "db_password", "secret", "Database password")
	flag.StringVar(&config.DbName, "db_name", "metadata_db", "Database name")
	flag.StringVar(&config.DbSslMode, "db_sslmode", "disable", "Database sll mode")
	flag.Parse() */
	config.Domain = os.Getenv("DOMAIN")
	config.Port = os.Getenv("PORT")
	config.ConsulAddress = os.Getenv("CONSUL_ADDRESS")
	config.KafkaAddress = os.Getenv("KAFKA_ADDRESS")
	config.DbDriver = os.Getenv("DB_DRIVER")
	config.DbConnection = os.Getenv("DB_CONNECTION")
	return &config
}
