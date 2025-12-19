package configs

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

var Config *config

type (
	config struct {
		App    app     `yaml:"app"`
		DB     db      `yaml:"db"`
		Redis  redis   `yaml:"redis"`
		Oauth  oauth   `yaml:"oauth"`
		Mail   mail    `yaml:"mail"`
	}

	app struct {
		Name        string `yaml:"name"`
		Url         string `yaml:"url"`
	}

	db struct {
		Postgres postgres `yaml:"postgres"`
	}

	postgres struct {
		Connection string `yaml:"connection"`
		Host       string `yaml:"host"`
		Port       int    `yaml:"port"`
		Database   string `yaml:"database"`
		Username   string `yaml:"username"`
		Password   string `yaml:"password"`
		SSLMode    string `yaml:"sslmode"`
	}

	redis struct {
		Host     string `yaml:"host"`
		Password string `yaml:"password"`
		Port     string `yaml:"port"`
	}

	oauth struct {
		AccessTokenExp  int64  `yaml:"access_token_exp"`
		RefreshTokenExp int64  `yaml:"refresh_token_exp"`
		ClientID        string `yaml:"client_id"`
		ClientSecret    string `yaml:"client_secret"`
	}
	mail struct {
		Host        string `yaml:"host"`
		Driver      string `yaml:"driver"`
		Port        int    `yaml:"port"`
		Username    string `yaml:"username"`
		Password    string `yaml:"password"`
		Encryption  string `yaml:"encryption"`
		FromAddress string `yaml:"from_address"`
	}
)

func InitConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("не удалось спарсить конфиг файл! Ошибка:%s", err)

		log.Fatal(err)
	}

	Config = &config{
		App: app{
			Name:        viper.GetString("app.name"),
			Url:         viper.GetString("app.url"),
		},
		DB: db{
			Postgres: postgres{
				Connection: viper.GetString("db.postgres.connection"),
				Host:       viper.GetString("db.postgres.host"),
				Port:       viper.GetInt("db.postgres.port"),
				Database:   viper.GetString("db.postgres.database"),
				Username:   viper.GetString("db.postgres.username"),
				Password:   viper.GetString("db.postgres.password"),
				SSLMode:    viper.GetString("db.postgres.sslmode"),
			},
		},
		Redis: redis{
			Host:     viper.GetString("redis.host"),
			Password: viper.GetString("redis.password"),
			Port:     viper.GetString("redis.port"),
		},
		Oauth: oauth{
			AccessTokenExp:  viper.GetInt64("oauth.access_token_exp"),
			RefreshTokenExp: viper.GetInt64("oauth.refresh_token_exp"),
			ClientID:        "kinokor-auth-client-id",
			ClientSecret:    "kinokor-auth-client-secret",
		},
		Mail: mail{
			Host:        viper.GetString("mail.host"),
			Driver:      viper.GetString("mail.driver"),
			Port:        viper.GetInt("mail.port"),
			Username:    viper.GetString("mail.username"),
			Password:    viper.GetString("mail.password"),
			Encryption:  viper.GetString("mail.encryption"),
			FromAddress: viper.GetString("mail.from_address"),
		},
	}
}
