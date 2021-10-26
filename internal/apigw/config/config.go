package config

import "github.com/spf13/viper"

type Config struct {
	Port            string `mapstructure:"PORT"`
	OrderGRPCurl    string `mapstructure:"ORDER_GRPC_URL"`
	OrderGRPCport   string `mapstructure:"ORDER_GRPC_PORT"`
	PaymentGRPCurl  string `mapstructure:"PAYMENT_GRPC_URL"`
	PaymentGRPCport string `mapstructure:"PAYMENT_GRPC_PORT"`
}

func NewConfig() (config *Config, err error) {
	viper.AddConfigPath("./config")
	viper.SetConfigName(".env.apigw")
	viper.SetConfigType("env")
	viper.AutomaticEnv()
	if err = viper.ReadInConfig(); err != nil {
		return
	}
	err = viper.Unmarshal(&config)
	return
}
