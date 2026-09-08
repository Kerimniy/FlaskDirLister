package tgbot

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type AdminStructTG struct {
	ChatID   int64  `json:"chat_id"`
	UserName string `json:"username"`
}

var AdminTG AdminStructTG = AdminStructTG{ChatID: -1}
var b *bot.Bot
var ctx context.Context

func InitTGBot() {

	InitChatID()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	opts := []bot.Option{
		bot.WithDefaultHandler(handler),
	}

	b, err := bot.New(os.Getenv("TELEGRAM_BOT_TOKEN"), opts...)
	if nil != err {

		panic(err)
	}

	b.Start(ctx)

}

func SaveChatID() {
	b, err := json.Marshal(AdminTG)

	if err != nil {
		log.Fatalln("could not serialize admin data")
	}

	f, err := os.Create(".TG_CHAT_ID")

	if err != nil {
		log.Fatalln("Could not write chat id file")
	}

	_, err = f.Write(b)
	if err != nil {
		log.Fatalln("Could not write chat id file")
	}
}

func InitChatID() {
	f, err := os.Open(".TG_CHAT_ID")
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return
		} else {
			log.Fatalln("Could not open chat id file")
		}

	}

	b, err := io.ReadAll(f)

	if err != nil {
		log.Fatalln("Could not read chat id file", err)
	}

	err = json.Unmarshal(b, &AdminTG)

	if err != nil {
		os.Remove(".TG_CHAT_ID")
		log.Fatalln("Invalid chat id")

	}

}

func SendCode(code string) error {

	kb := &models.ReplyKeyboardMarkup{
		Keyboard: [][]models.KeyboardButton{
			{
				{
					Text: "Get code",
				},
			},
		},
		ResizeKeyboard:  true,
		OneTimeKeyboard: false,
	}

	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      AdminTG.ChatID,
		Text:        code,
		ReplyMarkup: kb,
	})

	return err

}

func handler(ctx context.Context, b *bot.Bot, update *models.Update) {
	kb := &models.ReplyKeyboardMarkup{
		Keyboard: [][]models.KeyboardButton{
			{
				{
					Text: "Get code",
				},
			},
		},
		ResizeKeyboard:  true,
		OneTimeKeyboard: false,
	}

	if update.Message != nil {

		if AdminTG.ChatID == -1 || AdminTG.UserName == "" {

			AdminTG.ChatID = update.Message.Chat.ID
			AdminTG.UserName = update.Message.Chat.Username

			SaveChatID()

			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID:      update.Message.Chat.ID,
				Text:        fmt.Sprintf("Admin registered, chat ID: %d, username: %s", update.Message.Chat.ID, update.Message.Chat.Username),
				ReplyMarkup: kb,
			})
		} else if AdminTG.ChatID == update.Message.Chat.ID && AdminTG.UserName == update.Message.Chat.Username {
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID:      update.Message.Chat.ID,
				Text:        "No code for now",
				ReplyMarkup: kb,
			})
		}

	}
}
