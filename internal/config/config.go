package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
)

const defaultConfigFilename = "config.yml"

// Paths to search files in, on order of priority, latest takes precedence
var configPaths = []string{"/usr/local/etc/", "./configs/"}

// Config struct with field name
type Config struct {
	LogFile          string `yaml:"log_file" env:"YB_LOG_FILE" env-default:""`
	LogLevel         string `yaml:"log_level" env:"YB_LOG_LEVEL" env-default:"info"`
	TelegramBotToken string `yaml:"telegram_bot_token" env:"YB_TELEGRAM_BOT_TOKEN"`
}

func ReadConfig() *Config {
	var conf Config

	for _, file := range getConfigFiles() {
		_ = cleanenv.ReadConfig(file, &conf)
	}

	_ = cleanenv.ReadEnv(&conf)

	return &conf
}

func getConfigFiles() []string {
	var filename string

	if filename == "" {
		filename = os.Getenv("YB_CONFIG")
	}

	if filename != "" {
		return []string{filename}
	}

	files := make([]string, 0)
	for _, path := range configPaths {
		files = append(
			files,
			fmt.Sprintf("%s%s", path, defaultConfigFilename))
	}
	return files
}
