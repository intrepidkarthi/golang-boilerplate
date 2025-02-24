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
	Host            string        `mapstructure:"DB_HOST"`
	Port            int           `mapstructure:"DB_PORT"`
	User            string        `mapstructure:"DB_USER"`
	Password        string        `mapstructure:"DB_PASSWORD"`
	DBName          string        `mapstructure:"DB_NAME"`
	SSLMode         string        `mapstructure:"DB_SSLMODE"`
	MaxOpenConns    int           `mapstructure:"DB_MAX_OPEN_CONNS"`
	MaxIdleConns    int           `mapstructure:"DB_MAX_IDLE_CONNS"`
	ConnMaxLifetime time.Duration `mapstructure:"DB_CONN_MAX_LIFETIME"`
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
	// Server defaults
	viper.SetDefault("SERVER.PORT", "8080")
	viper.SetDefault("SERVER.READ_TIMEOUT", time.Second*10)
	viper.SetDefault("SERVER.WRITE_TIMEOUT", time.Second*10)
	viper.SetDefault("GRPC.PORT", "50051")

	// Database defaults
	viper.SetDefault("DB_PORT", 5432)
	viper.SetDefault("DB_MAX_OPEN_CONNS", 25)
	viper.SetDefault("DB_MAX_IDLE_CONNS", 5)
	viper.SetDefault("DB_CONN_MAX_LIFETIME", time.Minute*15)
	viper.SetDefault("DB_SSLMODE", "disable")
	
	viper.AutomaticEnv()
	
	var config Config
	err := viper.Unmarshal(&config)
	if err != nil {
		return nil, err
	}
	
	return &config, nil
}
