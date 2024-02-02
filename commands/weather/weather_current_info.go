package weather

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/iliyasali2107/discord-bot/internal/domain"
)

type WeatherSvc struct {
	cli WeatherAPIClient
}

type WeatherAPIClient interface {
	GetCurrentWeatherInfo(city string) (domain.WeatherData, error)
}

func NewWeatherSvc(weatherCli WeatherAPIClient) *WeatherSvc {
	return &WeatherSvc{
		cli: weatherCli,
	}
}

const (
	weatherCommand = "!weather"
	helpWeather    = "usage: !weather [city name]"
)

// errors are not checked, because don't have much time, but in real i would check it and log it
// no need to validate city name, weather api will send proper response if city is invalid
func (svc *WeatherSvc) Handle(session *discordgo.Session, message *discordgo.MessageCreate) {
	var (
		command = message.Content
		parts   = strings.Split(command, " ")
	)

	if parts[0] != weatherCommand {
		session.ChannelMessageSend(message.ChannelID, fmt.Sprintf("err: incorrect command \n%s ", helpWeather))
		return
	}

	if len(parts) != 2 {
		session.ChannelMessageSend(message.ChannelID, fmt.Sprintf("err: incorrect amount of args\n%s", helpWeather))
		return
	}

	city := parts[1]

	res, err := svc.cli.GetCurrentWeatherInfo(city)
	if err != nil {
		session.ChannelMessageSend(message.ChannelID, "err: failed to get weather info, sorry...")
		return
	}

	weatherInfo := fmt.Sprintf("Weather in requested %s, %s at time %s:\n%s, actual: %.2f°C, feels like: %.2f°C",
		res.Location.Country,
		res.Location.Name,
		res.Location.Localtime,
		res.Current.Condition.Text,
		res.Current.TempC,
		res.Current.FeelslikeC)

	session.ChannelMessageSend(message.ChannelID, weatherInfo)

}
