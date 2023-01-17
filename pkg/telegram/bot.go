package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/olegvolkov91/pocketer-bot/pkg/config"
	"github.com/olegvolkov91/pocketer-bot/pkg/repository"
	"github.com/zhashkevych/go-pocket-sdk"
)

type Bot struct {
	bot             *tgbotapi.BotAPI
	pocketClient    *pocket.Client
	tokenRepository repository.TokenRepository
	redirectURL     string

	messages config.Messages
}

func NewBot(
	bot *tgbotapi.BotAPI,
	pocketClient *pocket.Client,
	tokenRepository repository.TokenRepository,
	redirectURL string,
	messages config.Messages) *Bot {
	return &Bot{
		bot,
		pocketClient,
		tokenRepository,
		redirectURL,
		messages,
	}
}

func (b *Bot) Start() error {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := b.bot.GetUpdatesChan(u)
	if err != nil {
		return err
	}

	for update := range updates {
		if update.Message == nil {
			continue
		}

		if update.Message.IsCommand() {
			if err := b.handleCommand(update.Message); err != nil {
				b.handleError(update.Message.Chat.ID, err)
			}
			continue
		}

		if err := b.handleMessage(update.Message); err != nil {
			b.handleError(update.Message.Chat.ID, err)
		}
	}

	return nil
}
