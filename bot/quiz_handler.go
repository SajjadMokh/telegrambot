package bot

import (
	"BOT/games"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func SendQuizQuestion(
	bot *tgbotapi.BotAPI,
	chatID int64,
	userID int64,
) {

	index := games.GetQuestionIndex(userID)

	if index >= len(games.SajjadQuiz) {

		score := games.GetScore(userID)

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

		msg := tgbotapi.NewMessage(
			chatID,
			result,
		)

		bot.Send(msg)

		return
	}

	question := games.SajjadQuiz[index]

	var rows [][]tgbotapi.InlineKeyboardButton

	for _, option := range question.Options {

		button := tgbotapi.NewInlineKeyboardButtonData(
			option,
			option,
		)

		rows = append(
			rows,
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

	msg := tgbotapi.NewMessage(
		chatID,
		text,
	)

	msg.ReplyMarkup = keyboard

	bot.Send(msg)

}
