package config

import "github.com/spf13/viper"

type Config struct {
	Port            string `mapstructure:"PORT"`
	URL             string `mapstructure:"URL"`
	MongoURL        string `mapstructure:"MONGO_URL"`
	PaymentGRPCurl  string `mapstructure:"PAYMENT_GRPC_URL"`
	PaymentGRPCport string `mapstructure:"PAYMENT_GRPC_PORT"`
}

func NewConfig() (config *Config, err error) {
	viper.AddConfigPath("./config")
	viper.SetConfigName(".env.order")
	viper.SetConfigType("env")
	viper.AutomaticEnv()
	if err = viper.ReadInConfig(); err != nil {
		return
	}
	err = viper.Unmarshal(&config)
	return
}
