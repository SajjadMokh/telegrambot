package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"

	"BOT/bot"
	"BOT/database"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

func main() {

	// =============================
	// Load Environment Variables
	// =============================

	if err := godotenv.Load(); err != nil {
		log.Println(".env file not found, using environment variables")
	}

	// =============================
	// Telegram Bot Token
	// =============================

	token := os.Getenv("BOT_TOKEN")

	if token == "" {
		log.Fatal("BOT_TOKEN missing")
	}

	// =============================
	// Database
	// =============================

	databaseURL := os.Getenv("DATABASE_URL")

	if databaseURL == "" {
		log.Fatal("DATABASE_URL missing")
	}

	db, err := database.Connect(databaseURL)

	if err != nil {
		log.Fatal("Database connection error:", err)
	}

	defer db.Close()

	log.Println("Database connected successfully")

	// =============================
	// Telegram Bot
	// =============================

	tgBot, err := tgbotapi.NewBotAPI(token)

	if err != nil {
		log.Fatal("Telegram bot initialization error:", err)
	}

	log.Println("Bot started:", tgBot.Self.UserName)

	// =============================
	// Render URL
	// =============================

	domain := os.Getenv("RENDER_EXTERNAL_URL")

	if domain == "" {
		log.Fatal("RENDER_EXTERNAL_URL missing")
	}

	webhookURL := strings.TrimRight(domain, "/") + "/telegram"

	// =============================
	// Telegram Webhook
	// =============================

	webhook, err := tgbotapi.NewWebhook(webhookURL)

	if err != nil {
		log.Fatal("Create webhook error:", err)
	}

	_, err = tgBot.Request(webhook)

	if err != nil {
		log.Fatal("Set webhook error:", err)
	}

	info, err := tgBot.GetWebhookInfo()

	if err != nil {
		log.Fatal("Get webhook info error:", err)
	}

	log.Println("Webhook:", info.URL)

	// =============================
	// Telegram Update Endpoint
	// =============================

	http.HandleFunc("/telegram", func(w http.ResponseWriter, r *http.Request) {

		var update tgbotapi.Update

		if err := json.NewDecoder(r.Body).Decode(&update); err != nil {

			log.Println("Decode error:", err)

			w.WriteHeader(http.StatusBadRequest)

			return
		}

		// ارسال Database به Handler
		go bot.HandleUpdate(
			tgBot,
			update,
			db,
		)

		w.WriteHeader(http.StatusOK)
	})

	// =============================
	// Health Check
	// =============================

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		w.WriteHeader(http.StatusOK)

		_, _ = w.Write(
			[]byte("Telegram Bot Running"),
		)
	})

	// =============================
	// HTTP Server
	// =============================

	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	log.Println("Listening on port:", port)

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}