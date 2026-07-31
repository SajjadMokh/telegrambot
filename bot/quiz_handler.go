package bot

import (
	"BOT/games"
	"log"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func SendQuizQuestion(
	bot *tgbotapi.BotAPI,
	chatID int64,
	userID int64,
) {

	log.Println("========================================")
	log.Println("ENTER SendQuizQuestion")

	index := games.GetQuestionIndex(userID)

	log.Printf("UserID: %d", userID)
	log.Printf("Current Index: %d", index)

	if index >= len(games.SajjadQuiz) {

		log.Println("Quiz Finished")

		score := games.GetScore(userID)

		log.Printf("Final Score: %d", score)

		var result string

		switch score {

		case 5:

			result = `
🏆 فوق العاده!

امتیاز: 5/5

🔥 تو سجاد رو از خودش هم بهتر میشناسی!
رسماً عضو تیم هواداران ویژه شدی 😂👑
`

		case 4:

			result = `
🔥 عالی بود!

امتیاز: 4/5

تقریباً سجاد رو کامل شناختی 😎
فقط یک جا اشتباه زدی 😂
`

		case 3:

			result = `
😎 بد نبود!

امتیاز: 3/5

اطلاعاتت قابل قبوله ولی هنوز جای پیشرفت داری 😂
`

		case 2:

			result = `
😂 متوسط!

امتیاز: 2/5

فکر کنم باید بیشتر با سجاد وقت بگذرونی 🤣
`

		default:

			result = `
💀 فاجعه تاریخی!

امتیاز: 0 یا 1 از 5

ناموسا سجاد رو از نزدیک تا حالا دیدی اصلا 😂
برو شناخت بیشتر بگیر بعد بیا!
`

		}

		msg := tgbotapi.NewMessage(chatID, result)

		sent, err := bot.Send(msg)
		if err != nil {

			log.Println("Send Final Result Error:", err)

		} else {

			log.Printf("Final Result Sent. MessageID=%d", sent.MessageID)

		}

		log.Println("EXIT SendQuizQuestion")
		log.Println("========================================")

		return
	}

	question := games.SajjadQuiz[index]

	log.Printf("Sending Question #%d", index+1)
	log.Printf("Question Text: %s", question.Question)
	log.Printf("Correct Answer: %s", question.Answer)

	var rows [][]tgbotapi.InlineKeyboardButton

	for i, option := range question.Options {

		log.Printf("Option %d: %s", i+1, option)

		button := tgbotapi.NewInlineKeyboardButtonData(
			option,
			option,
		)

		rows = append(rows,
			[]tgbotapi.InlineKeyboardButton{
				button,
			},
		)

	}

	keyboard := tgbotapi.NewInlineKeyboardMarkup(rows...)

	text :=
		"❓ سوال " +
			strconv.Itoa(index+1) +
			" از 5\n\n" +
			question.Question

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ReplyMarkup = keyboard

	sent, err := bot.Send(msg)

	if err != nil {

		log.Println("Send Question Error:", err)

	} else {

		log.Printf("Question Sent Successfully. MessageID=%d", sent.MessageID)

	}

	log.Println("EXIT SendQuizQuestion")
	log.Println("========================================")
}