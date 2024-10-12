package cfg

import "github.com/spf13/viper"

func InitConfig() error {
	viper.AddConfigPath("cfg")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
