package config

import "flag"

type Config struct {
	Domain        string
	Port          string
	ConsulAddress string
	KafkaAddress  string
}

func NewConfig() *Config {
	config := Config{}
	flag.StringVar(&config.Domain, "domain", "localhost", "")
	flag.StringVar(&config.Port, "port", ":8082", "Listen for requests on this port")
	flag.StringVar(&config.ConsulAddress, "consul_address", "localhost:8500", "Register for service discovery")
	flag.StringVar(&config.KafkaAddress, "kafka_address", "localhost:9092", "Publish and consume messages")
	flag.Parse()
	return &config
}
