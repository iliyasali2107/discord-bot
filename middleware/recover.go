package middleware

import (
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/iliyasali2107/discord-bot/commands"
)

func recoverMiddleware(next commands.Command) commands.Command {
	return commands.HandlerFunc(func(session *discordgo.Session, message *discordgo.MessageCreate) {
		defer func() {
			err := recover()
			if err != nil {
				log.Printf("Recovered error: %v", err)
				session.ChannelMessageSend(message.ChannelID, "err: something happened to me and I couldn't process command")
				return
			}

		}()

		next.Handle(session, message)
	})
}

// apply recover to all commands

func ApplyRecover(g *commands.MessageCreateGateway) {
	g.WeatherSvc = recoverMiddleware(g.WeatherSvc)
	g.TranslatorSvc = recoverMiddleware(g.TranslatorSvc)
	g.GameSvc = recoverMiddleware(g.GameSvc)
	g.PollSvc = recoverMiddleware(g.PollSvc)
}
