package config

import (
	"bytes"
	"log"
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	Smtp     Smtp
	Database Database
	Grpc     Grpc
	Rabbit   Rabbit
	Gateway  Gateway
	Metrics  Metrics
}

type Smtp struct {
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

func LoadConfigFromEnv() (*Config, error) {
	viper := viper.New()
	cfgEnv := os.Getenv("CONFIG")
	var c Config

	if cfgEnv == "" {
		log.Fatal("Provide config env variable")
	}

	viper.AutomaticEnv()
	viper.SetConfigType("yaml")
	err := viper.ReadConfig(bytes.NewReader([]byte(cfgEnv)))

	if err != nil {
		log.Fatalf("Cannot read config %s", err.Error())
	}

	err = viper.Unmarshal(&c)

	log.Println(c)

	if err != nil {
		log.Fatal(err.Error())
	}

	return &c, nil
}
