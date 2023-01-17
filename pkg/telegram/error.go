package telegram

import (
	"errors"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var (
	errInvalidURL   = errors.New("url is invalid")
	errUnauthorized = errors.New("user is not authorized")
	errUnableToSave = errors.New("unable to save")
)

// msg.Text = "Ты не авторизирован! Используй команду /start"
// msg.Text = "Это невалидная ссылка!"
// msg.Text = "Увы, неудалось сохранить ссылку. Попробуй ещё раз позже"

func (b *Bot) handleError(chatID int64, err error) {
	msg := tgbotapi.NewMessage(chatID, "")

	switch err {
	case errInvalidURL:
		msg.Text = b.messages.InvalidURL
	case errUnauthorized:
		msg.Text = b.messages.Unauthorized
	case errUnableToSave:
		msg.Text = b.messages.UnableToSave
	default:
		msg.Text = b.messages.Default
	}
	b.bot.Send(msg)
}
