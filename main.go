package main

import (
	"log"
	"net/http"
	"os"

	"BOT/bot"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

func main() {

	godotenv.Load()

	token := os.Getenv("BOT_TOKEN")

	if token == "" {
		log.Fatal("BOT_TOKEN is missing")
	}

	go func() {

		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

			w.Write([]byte("Telegram Bot is running"))

		})

		port := os.Getenv("PORT")

		if port == "" {
			port = "8080"
		}

		http.ListenAndServe(":"+port, nil)

	}()

	tgBot, err := tgbotapi.NewBotAPI(token)

	if err != nil {
		log.Fatal(err)
	}

	log.Println("Bot started")

	u := tgbotapi.NewUpdate(0)

	u.Timeout = 60

	updates := tgBot.GetUpdatesChan(u)

	for update := range updates {

		bot.HandleUpdate(tgBot, update)

	}

}
