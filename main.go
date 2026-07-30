package main

import (
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}

	token := os.Getenv("BOT_TOKEN")

	if token == "" {
		log.Fatal("BOT_TOKEN is not set")
	}

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic("خطا در اتصال به تلگرام: ", err)
	}

	log.Printf("ربات %s با موفقیت روشن شد!", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {

		if update.Message == nil {
			continue
		}

		if update.Message.Text == "/start" {

			text := "سلام دوستان به ربات سجاد خوش آمدید"

			msg := tgbotapi.NewMessage(
				update.Message.Chat.ID,
				text,
			)

			_, err := bot.Send(msg)

			if err != nil {
				log.Println("خطا در ارسال پیام:", err)
			}
		}
	}
}