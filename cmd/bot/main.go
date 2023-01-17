package main

import (
	"log"

	"github.com/boltdb/bolt"
	"github.com/zhashkevych/go-pocket-sdk"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/olegvolkov91/pocketer-bot/pkg/config"
	"github.com/olegvolkov91/pocketer-bot/pkg/repository"
	"github.com/olegvolkov91/pocketer-bot/pkg/repository/boltdb"
	"github.com/olegvolkov91/pocketer-bot/pkg/server"
	"github.com/olegvolkov91/pocketer-bot/pkg/telegram"
)

func main() {
	cfg, err := config.Init()

	if err != nil {
		log.Fatal(err)
	}
	bot, err := tgbotapi.NewBotAPI(cfg.TelegramToken)
	if err != nil {
		log.Fatal(err)
	}

	bot.Debug = true
	pocketClient, err := pocket.NewClient(cfg.PocketConsumerKey)
	if err != nil {
		log.Fatal(err)
	}

	db, err := initDB(cfg)
	if err != nil {
		log.Fatal(err)
	}

	tokenRepository := boltdb.NewTokenRepository(db)

	telegramBot := telegram.NewBot(bot, pocketClient, tokenRepository, cfg.AuthServerURL, cfg.Messages)

	authorizationServer := server.NewAuthorizationServer(pocketClient, tokenRepository, cfg.TelegramBotURL)

	go func() {
		if err := telegramBot.Start(); err != nil {
			log.Fatal(err)
		}
	}()

	if err := authorizationServer.Start(); err != nil {
		log.Fatal(err)
	}
}

func initDB(cfg *config.Config) (*bolt.DB, error) {
	db, err := bolt.Open(cfg.DBPath, 0600, nil)

	if err != nil {
		log.Fatal(err)
	}

	if err := db.Batch(func(tx *bolt.Tx) error {
		if _, err := tx.CreateBucketIfNotExists([]byte(repository.AccessTokens)); err != nil {
			log.Fatal(err)
		}
		if _, err := tx.CreateBucketIfNotExists([]byte(repository.RequestTokens)); err != nil {
			log.Fatal(err)
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return db, nil
}
