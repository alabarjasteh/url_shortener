package config

import (
	"errors"
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Server ServerConfig
	Mysql  MysqlConfig
	Redis  RedisConfig
}

type ServerConfig struct {
	AppVersion string
	Port       string
}

type MysqlConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Dbname   string
	Driver   string
}

type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

func Load(filename string) (*viper.Viper, error) {
	v := viper.New()

	v.SetConfigName(filename)
	v.AddConfigPath(".")
	v.AutomaticEnv()
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil, errors.New("config file not found")
		}
		return nil, err
	}

	return v, nil
}

func ParseConfig(v *viper.Viper) (*Config, error) {
	var c Config

	err := v.Unmarshal(&c)
	if err != nil {
		log.Printf("unable to decode into struct, %v", err)
		return nil, err
	}

	return &c, nil
}
