package config

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/spf13/viper"
)

type Config struct {
	Responses
	Keyboards
}

type Keyboards struct {
	Main    tgbotapi.ReplyKeyboardMarkup
	BuyList tgbotapi.ReplyKeyboardMarkup
	Fridge  tgbotapi.ReplyKeyboardMarkup
}

type Responses struct {
	Start string `mapstructure:"start"`
}

var (
	mainKeyboard = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Добавить новый продукт"),
			tgbotapi.NewKeyboardButton("Добавить продукт в холодильник"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Открыть продукт"),
			tgbotapi.NewKeyboardButton("Изменить статус"),
			tgbotapi.NewKeyboardButton("Получить статистку"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Список продуктов в холодильнике"),
			tgbotapi.NewKeyboardButton("Список ранее используемых продуктов"),
		),
	)
)

func Init() (*Config, error) {
	if err := setUpViper(); err != nil {
		return nil, err
	}

	var cfg Config
	if err := unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func unmarshal(cfg *Config) error {
	if err := viper.Unmarshal(&cfg); err != nil {
		return err
	}

	if err := viper.UnmarshalKey("config.response", &cfg.Responses); err != nil {
		return err
	}

	// if err := viper.UnmarshalKey("messages.error", &cfg.Messages.Errors); err != nil {
	// 	return err
	// }

	return nil
}

func setUpViper() error {
	viper.AddConfigPath("var")
	viper.SetConfigName("msg")

	return viper.ReadInConfig()
}

func (c *Config) NewKeyboard() Keyboards {
	return Keyboards{
		Main: mainKeyboard,
	}
}
