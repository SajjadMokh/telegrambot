package bot

import (
	"BOT/friends"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func HandleUpdate(bot *tgbotapi.BotAPI, update tgbotapi.Update) {

	// اگر پیام خالی است، کاری نکن
	if update.Message == nil {
		return
	}

	// فقط به دستور start جواب بده
	if !update.Message.IsCommand() || update.Message.Command() != "start" {
		return
	}

	// دریافت مستقیم Username از تلگرام
	username := update.Message.From.UserName

	var text string

	// بررسی Username در لیست دوستان
	if message, ok := friends.Friends[username]; ok {
		text = message
	} else {
		text = "سلام دوست جدید 👋 هنوز اطلاعاتی از تو ندارم"
	}

	// ارسال پیام
	msg := tgbotapi.NewMessage(
		update.Message.Chat.ID,
		text,
	)

	bot.Send(msg)
}