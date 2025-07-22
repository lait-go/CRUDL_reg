package config

import (
	"log"
	"sub/internal/models"

	"github.com/spf13/viper"
)

func InitConf() models.Conf{
	viper.SetConfigFile("configs/config.yaml")

	viper.SetDefault("port", 8080)

	err := viper.ReadInConfig()
	if err != nil{
		log.Fatal("Ошибка при считывании конфига")
	}

	dbKey := KeyGath(models.Database{
			Host:     viper.GetString("database.host"),
			Port:     viper.GetString("database.port"),
			User:     viper.GetString("database.user"),
			Password: viper.GetString("database.password"),
			Name:     viper.GetString("database.name"),
		})

	return models.Conf{
		Port: viper.GetInt("port"),
		DbKey: dbKey,
	}
}

func KeyGath(conf models.Database) string{
	if conf.Host == "" || conf.Port == "" || conf.User == "" || conf.Password == "" || conf.Name == "" {
		return ""
	}

	return "postgres://" + conf.User + ":" + conf.Password + "@" + conf.Host + ":" + conf.Port + "/" + conf.Name + "?sslmode=disable"
}