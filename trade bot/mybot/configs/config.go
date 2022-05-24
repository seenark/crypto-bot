package configs

import (
	"strings"

	"github.com/spf13/viper"
)

type EchoConfig struct {
	Port int
}

type MongoConfig struct {
	Username       string
	Password       string
	DbName         string
	EmaCollection  string
	UserCollection string
}

type Configuration struct {
	Environment string
	EchoApp     EchoConfig
	Mongo       MongoConfig
}

func initConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.AddConfigPath("./configs")
	// read config from ENV
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	// read config
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}

func GetConfig() Configuration {
	initConfig()
	config := Configuration{}
	err := viper.Unmarshal(&config)
	if err != nil {
		panic(err)
	}
	return config
}
