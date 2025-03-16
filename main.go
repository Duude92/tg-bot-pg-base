package main

import (
	"context"
	"errors"
	"github.com/caarlos0/env/v11"
	"log"
	"os"
	"os/signal"
	"strings"
	Models "testTgPgBot/Models"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

// Send any text message to the bot after the bot has been started

type DbConfig struct {
	Host       string `env:"DB_HOST"`
	DBPort     int    `env:"DB_PORT"`
	DBUser     string `env:"DB_USER"`
	DBPassword string `env:"DB_PASSWORD"`
	DBName     string `env:"DB_NAME"`
}

func loadConfig() (*DbConfig, error) {
	var config DbConfig
	err := env.Parse(&config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func main() {
	if err := initDb(); err != nil {
		log.Fatalf("Unable to initialize database %v", err)
	}

	if err := initTgBot(); err != nil {
		log.Fatalf("Unable to initialize Telegram bot %v", err)
	}
}

func initTgBot() error {
	tgApiKey := os.Getenv("ApiKey")
	if tgApiKey == "" {
		return errors.New("Telegram Api key is not set")
	}
	tgApiKey = strings.TrimSpace(tgApiKey)
	opts := []bot.Option{
		bot.WithDefaultHandler(handler),
	}
	b, err := bot.New(tgApiKey, opts...)
	if err != nil {
		return err
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	b.RegisterHandler(bot.HandlerTypeMessageText, "/start", bot.MatchTypeExact, startHandler)
	b.RegisterHandler(bot.HandlerTypeMessageText, "/stop", bot.MatchTypeExact, stopHandler)
	b.RegisterHandler(bot.HandlerTypeMessageText, "/delete", bot.MatchTypeExact, deleteHandler)

	b.Start(ctx)
	return nil
}

func initDb() error {
	config, err := loadConfig()
	if err != nil {
		return err
	}
	return Models.CreateDb(config.Host, config.DBPort, config.DBUser, config.DBPassword, config.DBName)
}

func handler(ctx context.Context, b *bot.Bot, update *models.Update) {
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   update.Message.Text,
	})
}
func deleteHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	err := Models.DeleteUser(update.Message.Chat.ID)
	if err != nil {
		log.Fatalf("Unable to delete user %v", err)
	}
}
func startHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	err := Models.AddUser(&Models.User{
		Id:     update.Message.Chat.ID,
		Name:   update.Message.Chat.Username,
		Status: true,
	})
	if err != nil {
		log.Fatalf("Unable to add user %v", err)
	}
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   `User registered as ` + update.Message.Chat.Username,
	})
}
func stopHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	err := Models.UpdateUser(&Models.User{
		Id:     update.Message.Chat.ID,
		Name:   update.Message.Chat.Username,
		Status: false,
	})
	if err != nil {
		log.Fatalf("Unable to disable user %v", err)
	}
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   `User unregistered as ` + update.Message.Chat.Username,
	})
}
