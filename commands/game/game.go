package game

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
)

type GameSvc struct {
	secretNumber int
	isActive     bool
	mu           sync.Mutex
}

func NewGameSvc() *GameSvc {
	return &GameSvc{}
}

const (
	helpGame = "usage:\n - to start a game \" !game start\"\n- to make a guess \"!game guess [number]\""
)

func (svc *GameSvc) GameCommand(session *discordgo.Session, message *discordgo.MessageCreate) {
	parts := strings.Split(message.Content, " ")
	if parts[0] != "!game" {
		session.ChannelMessageSend(message.ChannelID, fmt.Sprintf("err: incorrect game input\n%s", helpGame))
		return
	}

	if parts[1] == "start" {
		if len(parts) != 2 {
			session.ChannelMessageSend(message.ChannelID, fmt.Sprintf("err: incorrect amount of start game args\n%s", helpGame))
			return
		}

		if svc.isActive {
			session.ChannelMessageSend(message.ChannelID, "err: game already started, finish this, then start new")
			return
		}

		go svc.start(session, message)

	} else if parts[1] == "guess" {
		if !svc.isActive {
			session.ChannelMessageSend(message.ChannelID, fmt.Sprintf("err: game is not active, start a game\n%s", helpGame))
			return
		}
		if len(parts) != 3 {
			session.ChannelMessageSend(message.ChannelID, fmt.Sprintf("err: incorrect amount of args\n%s", helpGame))
			return
		}

		num, err := strconv.Atoi(parts[2])
		if err != nil {
			session.ChannelMessageSend(message.ChannelID, fmt.Sprintf("err: guess should a number\n%s", helpGame))
			return
		}

		if num < 1 || num > 10 {
			session.ChannelMessageSend(message.ChannelID, fmt.Sprintf("err: guess should be between 1 to 10\n%s", helpGame))
			return
		}

		go svc.guess(session, message, num)
	} else {
		session.ChannelMessageSend(message.ChannelID, fmt.Sprintf("err: unknown command after !game \n%s", helpGame))
		return
	}
}

// locking and unlocking to prevent race condition,
// svc.secretNumber and svc.isActive are used in guess method

func (svc *GameSvc) start(session *discordgo.Session, message *discordgo.MessageCreate) {
	svc.mu.Lock()
	defer svc.mu.Unlock()
	rand.New(rand.NewSource(time.Now().UnixNano()))
	randomNumber := rand.Intn(5) + 1

	svc.secretNumber = randomNumber
	svc.isActive = true
	session.ChannelMessageSend(message.ChannelID, "game is started, guess the number between 1 and 5 inclusive\nto make a guess send \"!game guess [number]\"")
}

func (svc *GameSvc) guess(session *discordgo.Session, message *discordgo.MessageCreate, guessNum int) {
	svc.mu.Lock()
	defer svc.mu.Unlock()
	if guessNum == svc.secretNumber {
		svc.isActive = false
		session.ChannelMessageSend(message.ChannelID, fmt.Sprintf("yes, secret number was %d, and %s found it\n%s", svc.secretNumber, message.Author.Username, helpGame))
		return
	}

	session.ChannelMessageSend(message.ChannelID, "nooooooo, your guess is bad)")
}
