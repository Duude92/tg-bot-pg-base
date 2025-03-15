package main

import (
	"context"
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

func main() {
	tgApiKey := os.Getenv("ApiKey")
	//dbHost, dbPort, dbUser, dbPassword, dbName := os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWD"), os.Getenv("DB_NAME")
	var config DbConfig
	err := env.Parse(&config)
	if err != nil {
		log.Fatal(err)
	}
	if tgApiKey == "" {
		log.Fatal("Telegram Api key is not set")
	}
	tgApiKey = strings.TrimSpace(tgApiKey)
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	Models.CreateDb(config.Host, config.DBPort, config.DBUser, config.DBPassword, config.DBName)

	opts := []bot.Option{
		bot.WithDefaultHandler(handler),
	}

	b, err := bot.New(tgApiKey, opts...)
	if err != nil {
		panic(err)
	}
	b.RegisterHandler(bot.HandlerTypeMessageText, "/start", bot.MatchTypeExact, myStartHandler)
	b.RegisterHandler(bot.HandlerTypeMessageText, "/stop", bot.MatchTypeExact, myStopHandler)
	b.RegisterHandler(bot.HandlerTypeMessageText, "/delete", bot.MatchTypeExact, deleteHandler)

	b.Start(ctx)
}

func handler(ctx context.Context, b *bot.Bot, update *models.Update) {
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   update.Message.Text,
	})
}
func deleteHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	Models.DeleteUser(update.Message.Chat.ID)
}
func myStartHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	Models.AddUser(&Models.User{
		Id:     update.Message.Chat.ID,
		Name:   update.Message.Chat.Username,
		Status: true,
	})
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   `User registered as ` + update.Message.Chat.Username,
	})
}
func myStopHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	Models.UpdateUser(&Models.User{
		Id:     update.Message.Chat.ID,
		Name:   update.Message.Chat.Username,
		Status: false,
	})
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   `User unregistered as ` + update.Message.Chat.Username,
	})
}
