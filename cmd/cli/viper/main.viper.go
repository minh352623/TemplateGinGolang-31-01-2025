package main

import (
	"fmt"

	"github.com/spf13/viper"
)
type Config struct {
	Server struct {
		Port string `mapstructure:"port"`
		ReadTimeout string `mapstructure:"read_timeout"`
		WriteTimeout string `mapstructure:"write_timeout"`
	} `mapstructure:"server"`
	Databases []struct {
		Host string `mapstructure:"host"`
		Port string `mapstructure:"port"`
		User string `mapstructure:"user"`
		Password string `mapstructure:"password"`
		DBName string `mapstructure:"db_name"`
	} `mapstructure:"databases"`
}

func main() {
	viper := viper.New()
	viper.AddConfigPath("./config/")
	viper.SetConfigName("local")
	viper.SetConfigType("yaml")


	// config
	var config Config
	err := viper.Unmarshal(&config)
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	fmt.Println("config", config.Server.Port)

	for _, database := range config.Databases {
		fmt.Printf("database User: %s, Host: %s, Port: %s, DBName: %s\n", database.User, database.Host, database.Port, database.DBName)
	}
	
}
