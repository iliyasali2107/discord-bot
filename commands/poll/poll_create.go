package poll

import (
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

type PollSvc struct {
}

func NewPollSvc() *PollSvc {
	return &PollSvc{}
}

const (
	helpPoll = "usage: !poll -[question]- [options...]\n- question must be in -...-"

	// after pollDuration poll will end and show results
	pollDuration = 30 * time.Second
)

func (svc *PollSvc) Handle(session *discordgo.Session, message *discordgo.MessageCreate) {
	channelID := message.ChannelID

	// not standard input: question can have spaces
	// its better to use helper function
	res, err := validateAndParsePollCreateInput(message)
	if err != nil {
		if err == ErrInvalidCommand {
			session.ChannelMessageSend(channelID, fmt.Sprintf("err: incorrect !poll command \n%s ", helpPoll))
			return
		}

		if err == ErrInvalidAmountArgs {
			session.ChannelMessageSend(channelID, fmt.Sprintf("err: incorrect amount of args \n%s ", helpPoll))
			return
		}

		if err == ErrNoQuestion {
			session.ChannelMessageSend(channelID, fmt.Sprintf("err: no question provided \n%s ", helpPoll))
			return
		}
	}

	var (
		question    = res.question
		options     = res.options
		pollMessage = fmt.Sprintf("📊 ** %s  **\n duration is %v\n", question, pollDuration)
	)

	// add {number. option} to text, because emojis can't have text inside it
	for i, option := range options {
		pollMessage += fmt.Sprintf("%d. %s\n", i+1, strings.TrimSpace(option))
	}

	// send message
	msg, err := session.ChannelMessageSend(channelID, pollMessage)
	if err != nil {
		fmt.Println("Error sending poll message:", err)
		return
	}

	// add options as reactions
	for i := range options {
		err := session.MessageReactionAdd(channelID, msg.ID, fmt.Sprintf("%d⃣", i+1))
		if err != nil {
			fmt.Println("Error adding reaction:", err)
		}
	}

	// waits for poll ending
	time.Sleep(pollDuration)

	// get the message with reactions
	msg, err = session.ChannelMessage(channelID, msg.ID)
	if err != nil {
		fmt.Println("Error fetching message:", err)
		return
	}

	// process the reactions to determine the poll results
	results := make(map[string]int)
	for _, reaction := range msg.Reactions {
		results[reaction.Emoji.Name] = reaction.Count - 1
	}

	// display poll results
	resultMessage := "** [" + question + "] Poll Results:**\n\n"
	for i, option := range options {
		resultMessage += fmt.Sprintf("%s - %d votes\n", strings.TrimSpace(option), results[fmt.Sprintf("%d⃣", i+1)])
	}

	session.ChannelMessageSend(channelID, resultMessage)
}

type input struct {
	question string
	options  []string
}

var (
	ErrInvalidCommand    = fmt.Errorf("err: invalid !poll command")
	ErrInvalidAmountArgs = fmt.Errorf("err: incorrect amount of args")
	ErrNoQuestion        = fmt.Errorf("err: no question")
)

// i don't think its good validator, because it validates unstructured
func validateAndParsePollCreateInput(message *discordgo.MessageCreate) (input, error) {
	message.Content = strings.TrimSpace(message.Content)
	partsForValidation := strings.Split(message.Content, " ")
	if len(partsForValidation) < 3 {
		return input{}, ErrInvalidAmountArgs
	}

	if partsForValidation[0] != "!poll" {
		return input{}, ErrInvalidCommand
	}

	parts := strings.Split(message.Content, "-")
	if len(parts) == 2 {
		if parts[1] == message.Content {
			return input{}, ErrNoQuestion
		}

		return input{}, ErrNoQuestion
	}

	if len(parts) != 3 {
		return input{}, ErrNoQuestion
	}

	question := parts[1]
	options := strings.Split(strings.TrimSpace(parts[2]), " ")

	if len(options) < 1 || options[0] == "" {
		return input{}, ErrInvalidCommand
	}

	if len(question) == 0 {
		return input{}, ErrInvalidAmountArgs
	}

	return input{
		question: question,
		options:  options,
	}, nil

}
