package bot

import (
	"BOT/friends"
	"BOT/games"
	"log"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func HandleUpdate(bot *tgbotapi.BotAPI, update tgbotapi.Update) {

	// =============================
	// Callback Query (Inline Button)
	// =============================
	if update.CallbackQuery != nil {

		callback := update.CallbackQuery

		// جواب دادن سریع به تلگرام
		_, err := bot.Request(
			tgbotapi.NewCallback(
				callback.ID,
				"",
			),
		)

		if err != nil {
			log.Println("Callback answer error:", err)
		}

		if callback.Message == nil {

			log.Println("Callback message is nil")

			return
		}

		userID := callback.From.ID

		log.Println("========================================")
		log.Println("NEW CALLBACK")

		log.Printf(
			"UserID: %d Username: %s",
			userID,
			callback.From.UserName,
		)

		// دریافت شماره سوال و گزینه
		parts := strings.Split(
			callback.Data,
			"_",
		)

		if len(parts) != 2 {

			log.Println(
				"Invalid callback data:",
				callback.Data,
			)

			return
		}

		questionIndex, err := strconv.Atoi(parts[0])

		if err != nil {

			log.Println(
				"Question index error:",
				err,
			)

			return
		}

		optionIndex, err := strconv.Atoi(parts[1])

		if err != nil {

			log.Println(
				"Option index error:",
				err,
			)

			return
		}

		if questionIndex < 0 ||
			questionIndex >= len(games.SajjadQuiz) {

			log.Println(
				"Invalid question index:",
				questionIndex,
			)

			return
		}

		question := games.SajjadQuiz[questionIndex]

		if optionIndex < 0 ||
			optionIndex >= len(question.Options) {

			log.Println(
				"Invalid option index:",
				optionIndex,
			)

			return
		}

		answer := question.Options[optionIndex]
		log.Printf(
			"Clicked Answer: [%s]",
			answer,
		)

		currentIndex :=
			games.GetQuestionIndex(userID)

		log.Printf(
			"Current Question Index: %d",
			currentIndex,
		)

		if currentIndex >= len(games.SajjadQuiz) {

			log.Println(
				"Quiz already finished",
			)

			return
		}

		currentQuestion :=
			games.SajjadQuiz[currentIndex]

		log.Printf(
			"Question: %s",
			currentQuestion.Question,
		)

		log.Printf(
			"Correct Answer: [%s]",
			currentQuestion.Answer,
		)

		log.Printf(
			"Answer Match: %v",
			answer == currentQuestion.Answer,
		)

		if answer == currentQuestion.Answer {

			log.Println(
				"RESULT -> CORRECT",
			)

			games.AddScore(userID)

			log.Printf(
				"Score: %d",
				games.GetScore(userID),
			)

			msg := tgbotapi.NewMessage(
				callback.Message.Chat.ID,
				"✅ درست جواب دادی! +1 امتیاز 😎🔥",
			)

			_, err := bot.Send(msg)

			if err != nil {

				log.Println(
					"Send correct message error:",
					err,
				)

			}

		} else {

			log.Println(
				"RESULT -> WRONG",
			)

			msg := tgbotapi.NewMessage(
				callback.Message.Chat.ID,
				"❌ اشتباه بود 😂",
			)

			_, err := bot.Send(msg)

			if err != nil {

				log.Println(
					"Send wrong message error:",
					err,
				)
			}
		}

		log.Printf(
			"Index Before Next: %d",
			games.GetQuestionIndex(userID),
		)

		// حذف کامل پیام سوال قبلی
		deleteMsg := tgbotapi.NewDeleteMessage(
			callback.Message.Chat.ID,
			callback.Message.MessageID,
		)

		_, err = bot.Request(deleteMsg)

		if err != nil {
			log.Println("Delete Question Message Error:", err)
		}

		games.NextQuestion(userID)

		log.Printf(
			"Index After Next: %d",
			games.GetQuestionIndex(userID),
		)

		SendQuizQuestion(
			bot,
			callback.Message.Chat.ID,
			userID,
		)

		log.Println(
			"END CALLBACK",
		)

		log.Println(
			"========================================",
		)

		return
	}

	// =============================
	// Message Update
	// =============================

	if update.Message == nil {

		return
	}

	// =============================
	// Quiz Start
	// =============================

	if update.Message.IsCommand() &&
		update.Message.Command() == "quiz" {

		userID := update.Message.From.ID

		log.Printf(
			"QUIZ START UserID=%d",
			userID,
		)

		games.StartQuiz(
			userID,
		)

		SendQuizQuestion(
			bot,
			update.Message.Chat.ID,
			userID,
		)

		return
	}

	// =============================
	// Start Command
	// =============================

	if !update.Message.IsCommand() ||
		update.Message.Command() != "start" {

		return
	}

	username :=
		update.Message.From.UserName

	text :=
		"سلام دوست عزیز خوش اومدی 😎🔥"

	if message, ok := friends.Friends[username]; ok {

		text = message

	}

	msg := tgbotapi.NewMessage(
		update.Message.Chat.ID,
		text,
	)

	_, err := bot.Send(msg)

	if err != nil {

		log.Println(
			"Send start message error:",
			err,
		)

	}

	log.Printf(
		"New User Started Bot ID=%d Username=%s",
		update.Message.From.ID,
		username,
	)

}
