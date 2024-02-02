package clients

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/iliyasali2107/discord-bot/config"
	"github.com/iliyasali2107/discord-bot/internal/domain"
)

type WeatherApiV1 struct {
	url     string
	key     string
	current string
}

func NewWeatherApiV1(conf config.Config) *WeatherApiV1 {
	return &WeatherApiV1{
		url:     conf.WeatherApiURL,
		key:     conf.WeatherApiKey,
		current: conf.WeatherApiCurrent,
	}
}

func (api *WeatherApiV1) GetCurrentWeatherInfo(city string) (domain.WeatherData, error) {
	// Even it gets url from config, "key", "q", "aqi" are hardcoded(it is bad, in my opinion)
	url := fmt.Sprintf("%s/%s?key=%s&q=%s&aqi=no", api.url, api.current, api.key, city)
	response, err := http.Get(url)
	if err != nil {
		return domain.WeatherData{}, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return domain.WeatherData{}, fmt.Errorf("unsuccessful response, status code is %d", response.StatusCode)
	}

	// Decode the JSON response
	var weatherData domain.WeatherData
	err = json.NewDecoder(response.Body).Decode(&weatherData)
	if err != nil {
		return domain.WeatherData{}, err
	}

	return weatherData, nil
}
