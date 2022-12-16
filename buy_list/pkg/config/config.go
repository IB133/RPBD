package config

import (
	"github.com/go-co-op/gocron"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/spf13/viper"
)

type Config struct {
	Responses
	Keyboards
	Position
	Errors
	Scheduler
}

type Position struct {
	AddToBuyList           bool
	AddToFridge            bool
	OpenProduct            bool
	ChangeStatus           bool
	GetStatistic           bool
	UserInsert             bool
	AddToFridgeFromBuyList bool
}

type Scheduler struct {
	BuyListSched *gocron.Scheduler
	FridgeSched  *gocron.Scheduler
}

type Keyboards struct {
	Main     tgbotapi.ReplyKeyboardMarkup
	BuyOrNew tgbotapi.ReplyKeyboardMarkup
	Cancel   tgbotapi.ReplyKeyboardMarkup
	Current  tgbotapi.ReplyKeyboardMarkup
}

type Responses struct {
	Start               string `mapstructure:"start"`
	Succesfull          string `mapstructure:"succesfull"`
	AddingBuyList       string `mapstructure:"add_buy_list"`
	AddingFridgeList    string `mapstructure:"add_fridge_list"`
	AddingFridgeListBuy string `mapstructure:"add_fridge_list_from_buy"`
	AddingFridgeListNew string `mapstructure:"add_fridge_list_new"`
	GetStats            string `mapstructure:"get_stats"`
	StatusChange        string `mapstructure:"change_status"`
	ProductOpen         string `mapstructure:"open_product"`
}

type Errors struct {
	Default      string `mapstructure:"default"`
	UserNotFound string `mapstructure:"user_not_found"`
	ErrorInsert  string `mapstructure:"error_insert"`
	NoRows       string `mapstructure:"no_rows"`
}

var (
	mainKeyboard = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Добавить продукт в список покупок"),
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
	buylistOrNew = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Из списка покупок"),
			tgbotapi.NewKeyboardButton("Новый продукт"),
			tgbotapi.NewKeyboardButton("Отмена"),
		),
	)
	cancel = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Отмена"),
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

	if err := viper.UnmarshalKey("config.error", &cfg.Errors); err != nil {
		return err
	}

	return nil
}

func setUpViper() error {
	viper.AddConfigPath("var")
	viper.SetConfigName("msg")

	return viper.ReadInConfig()
}

func NewKeyboard() *Keyboards {
	return &Keyboards{
		Main:     mainKeyboard,
		BuyOrNew: buylistOrNew,
		Cancel:   cancel,
		Current:  mainKeyboard,
	}
}
