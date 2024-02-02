package translator

import (
	"context"
	"fmt"
	"strings"

	"cloud.google.com/go/translate"
	"github.com/bwmarrin/discordgo"
	iso6391 "github.com/emvi/iso-639-1"
	"golang.org/x/text/language"
)

type TranslatorSvc struct {
	cli *translate.Client
}

type Translator interface {
	Translate(targetLanguage, text string) (string, error)
}

func NewTranslator(cli *translate.Client) *TranslatorSvc {
	return &TranslatorSvc{
		cli: cli,
	}
}

const (
	translateCommand = "!translate"
	helpTranslate    = "usage: !translate [target language] [text to translate]"
)

func (svc *TranslatorSvc) TranslateCommand(session *discordgo.Session, message *discordgo.MessageCreate) {
	var (
		command   = message.Content
		parts     = strings.Split(command, " ")
		channelID = message.ChannelID
		ctx       = context.Background()
	)

	if parts[0] != translateCommand {
		session.ChannelMessageSend(channelID, fmt.Sprintf("err: incorrect command \n%s ", helpTranslate))
		return
	}

	if len(parts) < 3 {
		session.ChannelMessageSend(channelID, fmt.Sprintf("err: incorrect amount of args\n%s", helpTranslate))
		return
	}

	var (
		targetLanguage = parts[1]
		text           = strings.Join(parts[2:], " ")
	)

	// check that target language is real via iso6391 package, search by name, tag or native name, ex: {Russian, ru, Русский}
	// get iso6391.Language {Name, Code, NativeName}
	iso639Language, err := validateAndGetLanguage(targetLanguage)
	if err != nil {
		session.ChannelMessageSend(channelID, fmt.Sprintf("err: %s", err.Error()))
		return
	}

	// sending language to parser as tag to check if its correct tag, ex: en, ru
	langTag, err := language.Parse(iso639Language.Code)
	if err != nil {
		session.ChannelMessageSend(channelID, fmt.Sprintf("err: failed to parse target language \n%s", err.Error()))
		return
	}

	resp, err := svc.cli.Translate(ctx, []string{text}, langTag, nil)
	if err != nil {
		session.ChannelMessageSend(channelID, "err: failed to translate \n%s")
		return
	}

	if len(resp) == 0 {
		session.ChannelMessageSend(channelID, fmt.Sprintf("Translate returned empty response to text: %s", text))
		return
	}

	var (
		translation          = resp[0].Text
		sourceLanguageTag    = resp[0].Source.String()
		iso639SourceLanguage = iso6391.FromCode(sourceLanguageTag)
	)

	session.ChannelMessageSend(channelID, fmt.Sprintf("Translated from: %s\nTranslated to: %s\nTranslation: %s", iso639SourceLanguage.Name, iso639Language.Name, translation))

}

func capitalizeFirstLetter(s string) string {
	if s == "" {
		return s
	}

	firstLetter := strings.ToUpper(string(s[0]))
	return firstLetter + s[1:]
}

func validateAndGetLanguage(target string) (iso6391.Language, error) {
	lang := iso6391.Language{}
	target = strings.ToLower(target)

	// if valid code return iso6391.Language
	if iso6391.ValidCode(target) {
		return iso6391.FromCode(target), nil
	}

	// in iso6391 package's storage languages names and native names are stored capitalized
	target = capitalizeFirstLetter(target)
	lang = iso6391.FromName(target)
	if lang.Code != "" {
		return lang, nil
	}

	lang = iso6391.FromNativeName(target)
	if lang.Code == "" {
		return iso6391.Language{}, fmt.Errorf("failed to identify target language")
	}

	return lang, nil
}
