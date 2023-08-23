package config

import (
	"bytes"
	"os"
	"sync"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Config struct {
	SMTP     SMTP
	Database Database
	Grpc     Grpc
	Rabbit   Rabbit
	Gateway  Gateway
	Metrics  Metrics
	Jaeger   Jaeger
}

type SMTP struct {
	Host     string
	Port     int
	User     string
	Password string
}

type Database struct {
	Host     string
	User     string
	Password string
	DBName   string
	Port     int
	SslMode  string
}

type Grpc struct {
	Port int
}

type Gateway struct {
	Port int
}

type Metrics struct {
	Port int
}

type Rabbit struct {
	Host        string
	Port        int
	User        string
	Password    string
	QueueName   string
	ConsumePool int
}

type Jaeger struct {
	Host string
}

var config *Config

func LoadConfigFromEnv() *Config {
	viper := viper.New()
	cfgEnv := os.Getenv("CONFIG")

	var once sync.Once

	if cfgEnv == "" {
		log.WithFields(log.Fields{
			"message": "For running email service you need to provide CONFIG env variable, for more info see README.md",
		}).Fatal("Provider config env variable")
	}

	viper.AutomaticEnv()
	viper.SetConfigType("yaml")

	err := viper.ReadConfig(bytes.NewReader([]byte(cfgEnv)))

	if err != nil {
		log.WithFields(log.Fields{
			"message": "Unable to read yaml config file",
		}).Fatalf("Cannot read config %s", err.Error())
	}

	once.Do(func() {
		err = viper.Unmarshal(&config)
	})

	if err != nil {
		log.WithFields(log.Fields{
			"message": "Viper cannot unmarshal yaml config",
		}).Fatalf("Cannot unmarshal CONFIG variable yaml %s", err.Error())
	}

	return config
}
