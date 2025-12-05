package configs

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

var Config *config

type (
	config struct {
		App           app           `yaml:"app"`
		Redis         redis         `yaml:"redis"`
		Microservices microservices `yaml:"microservices"`
	}

	app struct {
		Name        string `yaml:"name"`
		Environment string `yaml:"environment"`
		Url         string `yaml:"url"`
		Key         string `yaml:"key"`
	}

	redis struct {
		Host     string `yaml:"host"`
		Password string `yaml:"password"`
		Port     string `yaml:"port"`
	}

	httpClient struct {
		Timeout       int `yaml:"timeout"`
		Retries       int `yaml:"retries"`
		BackoffMillis int `yaml:"backoff_millis"`
	}

	upstreamService struct {
		BaseURL string `yaml:"base_url"`
		ApiKey  string `yaml:"api_key"`
	}

	microservices struct {
		Http         httpClient      `yaml:"http"`
		Movies       upstreamService `yaml:"movies"`
		Auth         upstreamService `yaml:"auth"`
		Booking      upstreamService `yaml:"booking"`
		Notification upstreamService `yaml:"notification"`
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
		Redis: redis{
			Host:     viper.GetString("redis.host"),
			Password: viper.GetString("redis.password"),
			Port:     viper.GetString("redis.port"),
		},
		Microservices: microservices{
			Http: httpClient{
				Timeout:       viper.GetInt("microservices.http.timeout"),
				Retries:       viper.GetInt("microservices.http.retries"),
				BackoffMillis: viper.GetInt("microservices.http.backoff_millis"),
			},
			Movies: upstreamService{
				BaseURL: viper.GetString("microservices.movies.base_url"),
				ApiKey:  viper.GetString("microservices.movies.api_key"),
			},
			Auth: upstreamService{
				BaseURL: viper.GetString("microservices.auth.base_url"),
				ApiKey:  viper.GetString("microservices.auth.api_key"),
			},
			Notification: upstreamService{
				BaseURL: viper.GetString("microservices.notification.base_url"),
				ApiKey:  viper.GetString("microservices.notification.api_key"),
			},
			Booking: upstreamService{
				BaseURL: viper.GetString("microservices.booking.base_url"),
				ApiKey:  viper.GetString("microservices.booking.api_key"),
			},
		},
	}
}
