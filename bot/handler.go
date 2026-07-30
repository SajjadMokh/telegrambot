package bot

import (
	"BOT/friends"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func HandleUpdate(bot *tgbotapi.BotAPI, update tgbotapi.Update) {

	// 1. اگر پیام خالی است، هیچ کاری نکن
	if update.Message == nil {
		return
	}

	// 2. اگر پیام یک دستور نیست یا دستور آن "start" نیست، سکوت کن و پاسخ نده
	if !update.Message.IsCommand() || update.Message.Command() != "start" {
		return
	}

	// 3. دریافت یوزرنیم کاربر و تبدیل آن به حروف کوچک (برای جلوگیری از خطای بزرگ/کوچک بودن حروف)
	username := strings.ToLower(update.Message.From.UserName)

	var text string

	// 4. چک کردن یوزرنیم در مپ دوستان
	if message, ok := friends.Friends[username]; ok {
		text = message
	} else {
		text = "سلام دوست جدید 👋 هنوز اطلاعاتی از تو ندارم"
	}

	// 5. ارسال پیام
	msg := tgbotapi.NewMessage(
		update.Message.Chat.ID,
		text,
	)

	bot.Send(msg)
}
