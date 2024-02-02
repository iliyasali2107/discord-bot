package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	WeatherApiURL         string `yaml:"weather_api_url"`
	WeatherApiKey         string `yaml:"weather_api_key"`
	WeatherApiCurrent     string `yaml:"weather_api_current"`
	BotToken              string `yaml:"bot_token"`
	GoogleTranslateAPIKey string `yaml:"google_translate_api_key"`
}

func ParseYAML(fileName string) (Config, error) {
	body, err := os.ReadFile(fileName)
	if err != nil {
		return Config{}, err
	}
	settings := Config{}
	if err := yaml.Unmarshal(body, &settings); err != nil {
		return settings, err
	}
	return settings, nil
}
