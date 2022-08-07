package config

import (
	"github.com/spf13/viper"
	"log"
)

type Env struct {
	DBHostname    string
	DBPort        int
	DBUsername    string
	DBPassword    string
	DBDatabase    string
	DBMaxOpenConn int
	DBMaxIdleConn int
	HTTPServer    string
	HTTPPort      int
}

// use viper package to load/read the config file or .env file and
// return the value of the key
func LoadConfigFile() Env {
	var env Env

	viper.AddConfigPath("../config")
	viper.AddConfigPath("../../config")
	viper.SetConfigName("web-service")
	viper.SetConfigType("yaml")

	// Find and read the config file
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error while reading config file %s", err)
	}

	env.DBHostname, _ = viper.Get("clickhouse.host").(string)
	env.DBPort, _ = viper.Get("clickhouse.port").(int)
	env.DBUsername, _ = viper.Get("clickhouse.username").(string)
	env.DBPassword, _ = viper.Get("clickhouse.password").(string)
	env.DBDatabase, _ = viper.Get("clickhouse.database").(string)
	env.DBMaxOpenConn, _ = viper.Get("clickhouse.settings.maxOpenConn").(int)
	env.DBMaxIdleConn, _ = viper.Get("clickhouse.settings.maxIdleConn").(int)

	env.HTTPServer, _ = viper.Get("web-server.host").(string)
	env.HTTPPort, _ = viper.Get("web-server.port").(int)

	return env

}
