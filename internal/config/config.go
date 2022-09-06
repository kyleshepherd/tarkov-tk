package config

import (
	"github.com/rs/zerolog"
	"go.soon.build/kit/config"
)

// application name
const AppName = "tarkov-tk-bot"

// Config stores configuration options set by configuration file or env vars
type Config struct {
	Log      Log
	Discord  Discord
	Firebase Firebase
}

// Log contains logging configuration
type Log struct {
	Console bool
	Verbose bool
	Level   string
}

type Discord struct {
	BotToken       string
	RemoveCommands bool
	GuildID        string
}

type Firebase struct {
	ProjectID              string
	ServiceAccountFilePath string
}

// Default is a default configuration setup with sane defaults
var Default = Config{
	Log{
		Level: zerolog.InfoLevel.String(),
	},
	Discord{
		RemoveCommands: false,
		GuildID:        "",
	},
	Firebase{},
}

// New constructs a new Config instance
func New(opts ...config.Option) (Config, error) {
	c := Default
	v := config.ViperWithDefaults("tarkovtkbot")
	err := config.ReadInConfig(v, &c, opts...)
	if err != nil {
		return c, err
	}
	return c, nil
}
