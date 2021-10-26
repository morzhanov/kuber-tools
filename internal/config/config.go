package config

import "github.com/spf13/viper"

type Config struct {
	KafkaURL        string
	KafkaTopic      string
	KafkaGroupID    string
	MongoURL        string
	PostgresURL     string
	JaegerURL       string
	APIGWport       string
	OrderRESTurl    string
	PaymentGRPCurl  string
	PaymentGRPCport string
}

func NewConfig() (config *Config, err error) {
	viper.AddConfigPath("./config")
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AutomaticEnv()
	if err = viper.ReadInConfig(); err != nil {
		return
	}
	err = viper.Unmarshal(&config)
	return
}
