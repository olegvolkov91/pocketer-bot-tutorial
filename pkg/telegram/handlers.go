package telegram

import (
	"context"
	"net/url"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/zhashkevych/go-pocket-sdk"
)

const commandStart string = "start"

func (b *Bot) handleCommand(message *tgbotapi.Message) error {

	switch message.Command() {
	case commandStart:
		return b.handleStartCommand(message)
	default:
		return b.handleUnknownCommand(message)
	}
}

func (b *Bot) handleMessage(message *tgbotapi.Message) error {
	if _, err := url.ParseRequestURI(message.Text); err != nil {
		return errInvalidURL
	}

	accessToken, err := b.getAccessToken(message.Chat.ID)

	if err != nil {
		return errUnauthorized
	}

	if err := b.pocketClient.Add(context.Background(), pocket.AddInput{
		URL:         message.Text,
		AccessToken: accessToken,
	}); err != nil {
		return errUnableToSave
	}
	msg := tgbotapi.NewMessage(message.Chat.ID, b.messages.SavedSuccessfully)
	b.bot.Send(msg)
	return err
}

// func (b *Bot) saveLink(message *tgbotapi.Message, accessToken string) error {
// 	if err := b.validateURL(message.Text); err != nil {
// 		return err
// 	}

// 	if err := b.pocketClient.Add(context.Background(), pocket.AddInput{
// 		URL:         message.Text,
// 		AccessToken: accessToken,
// 	}); err != nil {
// 		return err
// 	}

// 	return nil
// }

func (b *Bot) handleStartCommand(message *tgbotapi.Message) error {
	_, err := b.getAccessToken(message.Chat.ID)
	if err != nil {
		return b.initAuthorizationProcess(message)
	}
	msg := tgbotapi.NewMessage(message.Chat.ID, b.messages.AlreadyAuthorized)
	_, err = b.bot.Send(msg)
	return err
}

func (b *Bot) handleUnknownCommand(message *tgbotapi.Message) error {
	msg := tgbotapi.NewMessage(message.Chat.ID, b.messages.UnknownCommand)
	_, err := b.bot.Send(msg)
	return err
}

func (b *Bot) validateURL(text string) error {
	_, err := url.ParseRequestURI(text)
	return err
}
