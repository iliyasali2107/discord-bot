package commands

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)

// create general interface to easily write recovere middleware.
// Because of recover works only in goroutine where panic occured
// i need to add middleware for all commands that works in goroutine
type Command interface {
	Handle(session *discordgo.Session, message *discordgo.MessageCreate)
}

func (hf HandlerFunc) Handle(session *discordgo.Session, message *discordgo.MessageCreate) {
	hf(session, message)
}

type HandlerFunc func(session *discordgo.Session, message *discordgo.MessageCreate)

type MessageCreateGateway struct {
	WeatherSvc    Command
	TranslatorSvc Command
	GameSvc       Command
	PollSvc       Command
}

func NewMessageCreateGateway(weatherSvc Command, translatorSvc Command, gameSvc Command, pollSvc Command) *MessageCreateGateway {
	return &MessageCreateGateway{
		WeatherSvc:    weatherSvc,
		TranslatorSvc: translatorSvc,
		GameSvc:       gameSvc,
		PollSvc:       pollSvc,
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
		go mainHelp(session, message)
	case strings.Contains(message.Content, "!weather"):
		go g.WeatherSvc.Handle(session, message)
	case strings.Contains(message.Content, "!translate"):
		go g.TranslatorSvc.Handle(session, message)
	case strings.Contains(message.Content, "!game"):
		go g.GameSvc.Handle(session, message)
	case strings.Contains(message.Content, "!poll"):
		go g.PollSvc.Handle(session, message)
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

func mainHelp(session *discordgo.Session, message *discordgo.MessageCreate) {
	session.ChannelMessageSend(message.ChannelID, help)
}
