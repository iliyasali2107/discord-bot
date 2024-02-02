package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"

	"cloud.google.com/go/translate"
	"github.com/bwmarrin/discordgo"
	"github.com/iliyasali2107/discord-bot/commands"
	"github.com/iliyasali2107/discord-bot/commands/game"
	"github.com/iliyasali2107/discord-bot/commands/poll"
	"github.com/iliyasali2107/discord-bot/commands/translator"
	"github.com/iliyasali2107/discord-bot/commands/weather"
	"github.com/iliyasali2107/discord-bot/config"
	"github.com/iliyasali2107/discord-bot/internal/clients"
	"github.com/iliyasali2107/discord-bot/middleware"
	"google.golang.org/api/option"
)

func main() {

	Run()
}

func Run() {
	var (
		conf     config.Config
		err      error
		confFlag bool
		botToken string

		// http://api.weatherapi.com/v1
		weatherAPIKey     string
		weatherAPIURL     string
		weatherAPICurrent string

		// official google cloud platform
		googleTranslateAPIKey string
	)

	flag.BoolVar(&confFlag, "config", false, "If true program starts using config")
	flag.StringVar(&botToken, "botToken", "", "Bot Token")
	flag.StringVar(&weatherAPIKey, "weatherAPIKey", "", "Weather API key from api.weatherapi.com/v1")
	flag.StringVar(&weatherAPIURL, "weatherAPIURL", "", "Weather API url")
	flag.StringVar(&weatherAPICurrent, "weatherAPICurrent", "current.json", "current.json for fetching current weather")
	flag.StringVar(&googleTranslateAPIKey, "googleTranslateAPIKey", "", "GCP API key")
	flag.Parse()

	if confFlag {
		conf, err = config.ParseYAML("./config.yaml")
		if err != nil {
			log.Fatalf("failed to parse config: %s", err.Error())
		}
	} else {
		conf.BotToken = botToken
		conf.WeatherApiKey = weatherAPIKey
		conf.WeatherApiURL = weatherAPIURL
		conf.WeatherApiCurrent = weatherAPICurrent
		conf.GoogleTranslateAPIKey = googleTranslateAPIKey
	}

	discord, err := discordgo.New("Bot " + conf.BotToken)
	if err != nil {
		log.Fatalf("failed to create discord session: %s\n", err.Error())
	}

	translatorCli, err := translate.NewClient(context.Background(), option.WithAPIKey(conf.GoogleTranslateAPIKey))
	if err != nil {
		log.Fatal(err)
	}

	var (
		translatorSvc        = translator.NewTranslator(translatorCli)
		weatherCli           = clients.NewWeatherApiV1(conf)
		weatherSvc           = weather.NewWeatherSvc(weatherCli)
		gameSvc              = game.NewGameSvc()
		polSvc               = poll.NewPollSvc()
		messageCreateGateway = commands.NewMessageCreateGateway(weatherSvc, translatorSvc, gameSvc, polSvc)
	)

	// appling recover middleware to all services
	middleware.ApplyRecover(messageCreateGateway)

	discord.AddHandler(messageCreateGateway.Handle)

	discord.Open()
	defer discord.Close()

	fmt.Println("Bot running...")

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	log.Println("shutting down")
}
