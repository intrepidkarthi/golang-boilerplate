package config

import (
	"github.com/spf13/viper"
	"time"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Redis    RedisConfig
	Kafka    KafkaConfig
	GRPC     GRPCConfig
}

type ServerConfig struct {
	Port         string        `mapstructure:"PORT"`
	ReadTimeout  time.Duration `mapstructure:"READ_TIMEOUT"`
	WriteTimeout time.Duration `mapstructure:"WRITE_TIMEOUT"`
}

type DatabaseConfig struct {
	Host     string `mapstructure:"DB_HOST"`
	Port     string `mapstructure:"DB_PORT"`
	User     string `mapstructure:"DB_USER"`
	Password string `mapstructure:"DB_PASSWORD"`
	DBName   string `mapstructure:"DB_NAME"`
	SSLMode  string `mapstructure:"DB_SSLMODE"`
}

type RedisConfig struct {
	Host     string `mapstructure:"REDIS_HOST"`
	Port     string `mapstructure:"REDIS_PORT"`
	Password string `mapstructure:"REDIS_PASSWORD"`
	DB       int    `mapstructure:"REDIS_DB"`
}

type KafkaConfig struct {
	Brokers []string `mapstructure:"KAFKA_BROKERS"`
	Topic   string   `mapstructure:"KAFKA_TOPIC"`
}

type GRPCConfig struct {
	Port string `mapstructure:"GRPC_PORT"`
}

func LoadConfig() (*Config, error) {
	viper.SetDefault("SERVER.PORT", "8080")
	viper.SetDefault("SERVER.READ_TIMEOUT", time.Second*10)
	viper.SetDefault("SERVER.WRITE_TIMEOUT", time.Second*10)
	viper.SetDefault("GRPC.PORT", "50051")
	
	viper.AutomaticEnv()
	
	var config Config
	err := viper.Unmarshal(&config)
	if err != nil {
		return nil, err
	}
	
	return &config, nil
}
