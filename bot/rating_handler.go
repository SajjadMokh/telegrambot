package bot

import (
	"database/sql"
	"log"
	"strconv"
	"strings"

	"BOT/service"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// ============================================================
// Handle Rating Callback
// ============================================================

func HandleRatingCallback(
	bot *tgbotapi.BotAPI,
	callback *tgbotapi.CallbackQuery,
	db *sql.DB,
) {
	if callback == nil || callback.Message == nil {
		return
	}

	userID := callback.From.ID
	chatID := callback.Message.Chat.ID
	data := callback.Data

	// ========================================================
	// Start Rating
	// ========================================================

	if strings.HasPrefix(data, "rate_") {

		teacherID, err := parseID(data, "rate_")
		if err != nil {
			log.Println("Invalid teacher ID:", err)
			return
		}

		msg := tgbotapi.NewMessage(
			chatID,
			"⭐ امتیازت به این استاد رو از ۱ تا ۱۰ انتخاب کن:",
		)

		msg.ReplyMarkup = ratingKeyboard(teacherID)

		if _, err := bot.Send(msg); err != nil {
			log.Println("Send rating keyboard error:", err)
		}

		return
	}

	// ========================================================
	// Rating Selected
	// ========================================================

	if strings.HasPrefix(data, "rating_") {

		parts := strings.Split(data, "_")

		if len(parts) != 3 {
			log.Println("Invalid rating callback:", data)
			return
		}

		// ====================================================
		// Teacher ID
		// ====================================================

		teacherID, err := strconv.ParseInt(
			parts[1],
			10,
			64,
		)

		if err != nil {
			log.Println("Invalid teacher ID:", err)
			return
		}

		// ====================================================
		// Rating
		// ====================================================

		rating, err := strconv.Atoi(parts[2])
		if err != nil {
			log.Println("Invalid rating:", err)
			return
		}

		if rating < 1 || rating > 10 {
			log.Println("Rating out of range:", rating)
			return
		}

		// ====================================================
		// Get Telegram Username
		// ====================================================

		username := ""

		if callback.From.UserName != "" {
			username = callback.From.UserName
		}

		// ====================================================
		// Rating Service
		// ====================================================

		ratingService := service.NewRatingService(db)

		// ====================================================
		// Get / Create User
		// ====================================================

		dbUserID, err := ratingService.GetOrCreateUser(
			userID,
			username,
		)

		if err != nil {
			log.Println("Get or create user error:", err)

			sendMessage(
				bot,
				chatID,
				"❌ خطایی در شناسایی کاربر اتفاق افتاد.",
			)

			return
		}

		// ====================================================
		// Save Rating
		// ====================================================

		err = ratingService.SaveRating(
			dbUserID,
			teacherID,
			rating,
		)

		if err != nil {
			log.Println("Save rating error:", err)

			sendMessage(
				bot,
				chatID,
				"❌ هنگام ثبت امتیاز مشکلی پیش اومد.",
			)

			return
		}

		// ====================================================
		// Success
		// ====================================================

		sendMessage(
			bot,
			chatID,
			"✅ امتیاز شما با موفقیت ثبت شد.\n\n"+
				"⭐ امتیاز: "+
				strconv.Itoa(rating)+
				" از ۱۰",
		)

		return
	}
}

// ============================================================
// Rating Keyboard
// ============================================================

func ratingKeyboard(
	teacherID int64,
) tgbotapi.InlineKeyboardMarkup {

	var rows [][]tgbotapi.InlineKeyboardButton
	var row []tgbotapi.InlineKeyboardButton

	for rating := 1; rating <= 10; rating++ {

		button := tgbotapi.NewInlineKeyboardButtonData(
			strconv.Itoa(rating)+" ⭐",
			"rating_"+
				strconv.FormatInt(
					teacherID,
					10,
				)+
				"_"+strconv.Itoa(rating),
		)

		row = append(row, button)

		// هر 5 امتیاز یک ردیف
		if len(row) == 5 {
			rows = append(rows, row)
			row = nil
		}
	}

	if len(row) > 0 {
		rows = append(rows, row)
	}

	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}