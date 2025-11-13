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
		Logger grayLog `yaml:"graylog"`
		Oauth  oauth   `yaml:"oauth"`
	}

	app struct {
		Name        string `yaml:"name"`
		Environment string `yaml:"environment"`
		Url         string `yaml:"url"`
		Key         string `yaml:"key"`
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

	grayLog struct {
		Host   string `yaml:"host"`
		Port   int    `yaml:"port"`
		Source string `yaml:"source"`
	}

	oauth struct {
		AccessTokenExp  int64  `yaml:"access_token_exp"`
		RefreshTokenExp int64  `yaml:"refresh_token_exp"`
		ClientID        string `yaml:"client_id"`
		ClientSecret    string `yaml:"client_secret"`
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
			Environment: viper.GetString("app.environment"),
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
		Logger: grayLog{
			Host:   viper.GetString("graylog.host"),
			Port:   viper.GetInt("graylog.port"),
			Source: viper.GetString("graylog.source"),
		},
		Oauth: oauth{
			AccessTokenExp:  viper.GetInt64("oauth.access_token_exp"),
			RefreshTokenExp: viper.GetInt64("oauth.refresh_token_exp"),
			ClientID:        "do-auth-client-id",
			ClientSecret:    "do-auth-client-secret",
		},
	}
}
