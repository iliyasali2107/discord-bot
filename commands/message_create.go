package commands

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)

type MessageCreateGateway struct {
	weatherSvc    WeatherInterface
	translatorSvc TranslatorInterface
	gameSvc       GameInterface
	pollSvc       PollInterface
}

// word Interface in naming in my opinion is not so good, but there are not much options
// Use interfaces because we can change services but not abstractions, and we doesn't depend on services
type WeatherInterface interface {
	WeatherCurrentInfoCommand(session *discordgo.Session, message *discordgo.MessageCreate)
}

type TranslatorInterface interface {
	TranslateCommand(session *discordgo.Session, message *discordgo.MessageCreate)
}

type GameInterface interface {
	GameCommand(session *discordgo.Session, message *discordgo.MessageCreate)
}

type PollInterface interface {
	CreatePoll(session *discordgo.Session, message *discordgo.MessageCreate)
}

func NewMessageCreateGateway(weatherSvc WeatherInterface, translatorSvc TranslatorInterface, gameSvc GameInterface, pollSvc PollInterface) *MessageCreateGateway {
	return &MessageCreateGateway{
		weatherSvc:    weatherSvc,
		translatorSvc: translatorSvc,
		gameSvc:       gameSvc,
		pollSvc:       pollSvc,
	}
}

func (g *MessageCreateGateway) Handle(session *discordgo.Session, message *discordgo.MessageCreate) {
	// prevent bot responding to its own message
	if message.Author.ID == session.State.User.ID {
		return
	}

	message.Content = strings.TrimSpace(message.Content)

	switch {
	case strings.Contains(message.Content, "!help"):
		go mainHelp(session, message.Message)
	case strings.Contains(message.Content, "!weather"):
		go g.weatherSvc.WeatherCurrentInfoCommand(session, message)
	case strings.Contains(message.Content, "!translate"):
		go g.translatorSvc.TranslateCommand(session, message)
	case strings.Contains(message.Content, "!game"):
		go g.gameSvc.GameCommand(session, message)
	case strings.Contains(message.Content, "!poll"):
		go g.pollSvc.CreatePoll(session, message)
	}
}

const help = `
bot usage
- !help - for list of commands 
- !weather [city]  - gets current weather info in city
- !translate [target language] [text]   -  translate text to target language
   target language can be in any format {"en", "English", "english"}
- !game start - to start "guess game" 
- !game guess [guess number] - to guess the secret number
- !poll -[question]- [options...]  - question must be in -...-
`

func mainHelp(session *discordgo.Session, message *discordgo.Message) {
	session.ChannelMessageSend(message.ChannelID, help)
}
