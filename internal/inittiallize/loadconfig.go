package inittiallize

import (
	"ecom/global"
	"fmt"
	"os"

	_ "github.com/joho/godotenv/autoload"
	"github.com/spf13/viper"
)

func LoadConfig() {
	environment := os.Getenv("APP_ENV")
	fmt.Println("environment", environment)
	if environment == "" {
		environment = "local"
	}

	viper := viper.New()
	viper.AddConfigPath("./config/")
	viper.SetConfigName(environment)
	viper.SetConfigType("yaml")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
	for _, key := range viper.AllKeys() {
		value := viper.GetString(key)
		if value != "" {
			viper.Set(key, os.ExpandEnv(value)) // os.ExpandEnv replaces ${VAR} with the value of VAR
		}
	}

	// config
	if err := viper.Unmarshal(&global.Config); err != nil {
		fmt.Println("error unmarshal config", err)
	}

	// for _, database := range global.Config.Databases {
	// 	fmt.Printf("database User: %s, Host: %s, Port: %s, DBName: %s\n", database.User, database.Host, database.Port, database.DBName)
	// }
}
