package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"

	"BOT/bot"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

func main() {

	godotenv.Load()

	token := os.Getenv("BOT_TOKEN")

	if token == "" {

		log.Fatal("BOT_TOKEN missing")

	}


	tgBot, err := tgbotapi.NewBotAPI(token)

	if err != nil {

		log.Fatal(err)

	}


	log.Println(
		"Bot started:",
		tgBot.Self.UserName,
	)



	domain := os.Getenv("RENDER_EXTERNAL_URL")

	if domain == "" {

		log.Fatal(
			"RENDER_EXTERNAL_URL missing",
		)

	}


	webhookURL := strings.TrimRight(
		domain,
		"/",
	) + "/telegram"



	webhook, err := tgbotapi.NewWebhook(
		webhookURL,
	)

	if err != nil {

		log.Fatal(
			"Create webhook error:",
			err,
		)

	}



	_, err = tgBot.Request(
		webhook,
	)


	if err != nil {

		log.Fatal(
			"Set webhook error:",
			err,
		)

	}



	info, err := tgBot.GetWebhookInfo()

	if err != nil {

		log.Fatal(err)

	}


	log.Println(
		"Webhook:",
		info.URL,
	)



	http.HandleFunc(
		"/telegram",
		func(
			w http.ResponseWriter,
			r *http.Request,
		) {


			var update tgbotapi.Update


			err := json.NewDecoder(
				r.Body,
			).Decode(
				&update,
			)


			if err != nil {

				log.Println(
					"Decode error:",
					err,
				)

				w.WriteHeader(
					http.StatusBadRequest,
				)

				return

			}



			go bot.HandleUpdate(
				tgBot,
				update,
			)



			w.WriteHeader(
				http.StatusOK,
			)

		},
	)



	http.HandleFunc(
		"/",
		func(
			w http.ResponseWriter,
			r *http.Request,
		){

			w.Write(
				[]byte(
					"Telegram Bot Running",
				),
			)

		},
	)



	port := os.Getenv("PORT")


	if port == "" {

		port = "8080"

	}



	log.Println(
		"Listening on port:",
		port,
	)



	err = http.ListenAndServe(
		":"+port,
		nil,
	)


	if err != nil {

		log.Fatal(err)

	}

}