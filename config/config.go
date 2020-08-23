package config

import "github.com/spf13/viper"

// See
type Config struct {
	Constants
}

type Constants struct {
	Line
	Spreadsheets
	Backend
}

// Line message API's secret and token
type Line struct {
	ChannelSecret string
	ChannelToken  string
}

type Spreadsheets struct {
	SpreadsheetId string
}

type Backend struct {
	CreateChatroomToken string
}

func New() (*Config, error) {
	config := Config{}
	constants, err := initViper()
	config.Constants = constants

	return &config, err
}

func initViper() (Constants, error) {
	viper.AddConfigPath(".")
	viper.SetConfigName("config")

	err := viper.ReadInConfig()
	if err != nil {
		return Constants{}, err
	}

	var constants Constants
	err = viper.Unmarshal(&constants)

	return constants, err
}
