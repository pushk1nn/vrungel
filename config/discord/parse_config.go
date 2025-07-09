package discord

import (
	"github.com/spf13/viper"
)

var (
	Name string

	Token string

	AppID string

	GuildID string

	Environment string
)

func init() {
	viper.SetConfigFile("config/discord/config.yml")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	UpdateConfigs()
}

func UpdateConfigs() {
	Name = viper.GetString("name")
	Token = viper.GetString("token")
	AppID = viper.GetString("app_id")
	GuildID = viper.GetString("guild_id")
	Environment = viper.GetString("env")
}
