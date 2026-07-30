package bot

import (
	"BOT/friends"
	"BOT/games"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func HandleUpdate(bot *tgbotapi.BotAPI, update tgbotapi.Update) {

	// =============================
	// کلیک روی دکمه شیشه‌ای
	// =============================
	if update.CallbackQuery != nil {

		callback := update.CallbackQuery

		userID := callback.From.ID

		answer := callback.Data

		log.Printf(
			"Quiz Answer | UserID: %d | Answer: %s",
			userID,
			answer,
		)

		index := games.GetQuestionIndex(userID)

		if index < len(games.SajjadQuiz) {

			question := games.SajjadQuiz[index]

			if answer == question.Answer {

				games.AddScore(userID)
				log.Printf(
					"QUIZ SCORE | UserID: %d | Username: %s | Score: %d",
					userID,
					callback.From.UserName,
					games.GetScore(userID),
				)

				msg := tgbotapi.NewMessage(
					callback.Message.Chat.ID,
					"✅ درست جواب دادی! +1 امتیاز 😎🔥",
				)

				_, err := bot.Send(msg)

				if err != nil {

					log.Println(
						"Send message error:",
						err,
					)

				}

			} else {

				msg := tgbotapi.NewMessage(
					callback.Message.Chat.ID,
					"❌ اشتباه بود 😂",
				)

				_, err := bot.Send(msg)

				if err != nil {

					log.Println(
						"Send message error:",
						err,
					)

				}

			}

			// رفتن به سوال بعدی
			games.NextQuestion(userID)

			SendQuizQuestion(
				bot,
				callback.Message.Chat.ID,
				userID,
			)

		}

		// بستن حالت loading روی دکمه تلگرام
		_, err := bot.Request(
			tgbotapi.NewCallback(
				callback.ID,
				"",
			),
		)

		if err != nil {

			log.Println(
				"Callback error:",
				err,
			)

		}

		return
	}

	// =============================
	// پیام معمولی
	// =============================

	if update.Message == nil {

		return
	}

	// =============================
	// شروع مسابقه
	// =============================

	if update.Message.IsCommand() &&
		update.Message.Command() == "quiz" {

		userID := update.Message.From.ID

		games.StartQuiz(userID)

		SendQuizQuestion(
			bot,
			update.Message.Chat.ID,
			userID,
		)

		return
	}

	// =============================
	// Start Bot
	// =============================

	if !update.Message.IsCommand() ||
		update.Message.Command() != "start" {

		return
	}

	log.Printf(
		"New User Started Bot | ID: %d | Username: %s | FirstName: %s",
		update.Message.From.ID,
		update.Message.From.UserName,
		update.Message.From.FirstName,
	)

	username := update.Message.From.UserName

	var text string

	if message, ok := friends.Friends[username]; ok {

		text = message

	} else {

		text = "سلام دوست عزیز خوش اومدی میتونی از بخش منو ها  سجاد شناسی رو انتخاب کنی و خودت رو بسنجی😎🔥.. "

	}

	msg := tgbotapi.NewMessage(
		update.Message.Chat.ID,
		text,
	)

	_, err := bot.Send(msg)

	if err != nil {

		log.Println(
			"Send message error:",
			err,
		)

	}

}
