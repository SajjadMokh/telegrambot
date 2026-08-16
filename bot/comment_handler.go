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
// Comment State
// ============================================================

type CommentState struct {
	TeacherID   int64
	IsAnonymous bool
	Step        int
}

var commentStates = make(map[int64]*CommentState)

// ============================================================
// Handle Comment Callback
// ============================================================

func HandleCommentCallback(
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
	// Start Comment
	// ========================================================

	if strings.HasPrefix(data, "comment_yes_") {

		teacherID, err := parseID(
			data,
			"comment_yes_",
		)

		if err != nil {
			log.Println("Invalid teacher ID:", err)
			return
		}

		commentStates[userID] = &CommentState{
			TeacherID: teacherID,
			Step:      1,
		}

		msg := tgbotapi.NewMessage(
			chatID,
			"💬 نظرت رو چطور ثبت کنم؟",
		)

		msg.ReplyMarkup = commentIdentityKeyboard(
			teacherID,
		)

		if _, err := bot.Send(msg); err != nil {
			log.Println(
				"Send comment identity keyboard error:",
				err,
			)
		}

		return
	}

	// ========================================================
	// Comment No
	// ========================================================

	if strings.HasPrefix(data, "comment_no_") {

		delete(
			commentStates,
			userID,
		)

		sendMessage(
			bot,
			chatID,
			"👍 باشه، فقط امتیازت ثبت شد.",
		)

		return
	}

	// ========================================================
	// Anonymous Comment
	// ========================================================

	if strings.HasPrefix(data, "comment_anonymous_") {

		teacherID, err := parseID(
			data,
			"comment_anonymous_",
		)

		if err != nil {
			log.Println("Invalid teacher ID:", err)
			return
		}

		commentStates[userID] = &CommentState{
			TeacherID:   teacherID,
			IsAnonymous: true,
			Step:        2,
		}

		sendMessage(
			bot,
			chatID,
			"✍️ حالا نظرت درباره استاد رو بنویس:\n\n"+
				"🕵️ نظر شما به صورت ناشناس ثبت میشه.",
		)

		return
	}

	// ========================================================
	// Public Comment
	// ========================================================

	if strings.HasPrefix(data, "comment_public_") {

		teacherID, err := parseID(
			data,
			"comment_public_",
		)

		if err != nil {
			log.Println("Invalid teacher ID:", err)
			return
		}

		commentStates[userID] = &CommentState{
			TeacherID:   teacherID,
			IsAnonymous: false,
			Step:        2,
		}

		sendMessage(
			bot,
			chatID,
			"✍️ حالا نظرت درباره استاد رو بنویس:\n\n"+
				"👤 نظر شما با آیدی تلگرام شما ثبت میشه.",
		)

		return
	}
}

// ============================================================
// Handle Comment Message
// ============================================================

func HandleCommentMessage(
	bot *tgbotapi.BotAPI,
	message *tgbotapi.Message,
	db *sql.DB,
) bool {

	if message == nil || message.From == nil {
		return false
	}

	userID := message.From.ID

	state, exists := commentStates[userID]

	if !exists {
		return false
	}

	if state.Step != 2 {
		return false
	}

	// ========================================================
	// Get Comment Text
	// ========================================================

	text := strings.TrimSpace(
		message.Text,
	)

	if text == "" {

		sendMessage(
			bot,
			message.Chat.ID,
			"❌ متن نظر نمی‌تونه خالی باشه.",
		)

		return true
	}

	// ========================================================
	// Get Telegram Username
	// ========================================================

	username := ""

	if message.From.UserName != "" {
		username = message.From.UserName
	}

	// ========================================================
	// Get / Create User
	// ========================================================

	ratingService := service.NewRatingService(db)

	dbUserID, err := ratingService.GetOrCreateUser(
		userID,
		username,
	)

	if err != nil {

		log.Println(
			"Get or create user error:",
			err,
		)

		sendMessage(
			bot,
			message.Chat.ID,
			"❌ خطایی هنگام شناسایی کاربر اتفاق افتاد.",
		)

		return true
	}

	// ========================================================
	// Save Comment
	// ========================================================

	commentService := service.NewCommentService(db)

	err = commentService.AddComment(
		state.TeacherID,
		dbUserID,
		text,
		state.IsAnonymous,
	)

	if err != nil {

		log.Println(
			"Add comment error:",
			err,
		)

		sendMessage(
			bot,
			message.Chat.ID,
			"❌ هنگام ثبت نظر مشکلی پیش اومد.",
		)

		return true
	}

	// ========================================================
	// Success
	// ========================================================

	sendMessage(
		bot,
		message.Chat.ID,
		"✅ نظرت با موفقیت ثبت شد.\n\n"+
			"ممنون که تجربه‌ات رو به اشتراک گذاشتی ❤️",
	)

	// ========================================================
	// Clear State
	// ========================================================

	delete(
		commentStates,
		userID,
	)

	return true
}

// ============================================================
// Comment Identity Keyboard
// ============================================================

func commentIdentityKeyboard(
	teacherID int64,
) tgbotapi.InlineKeyboardMarkup {

	return tgbotapi.NewInlineKeyboardMarkup(

		[]tgbotapi.InlineKeyboardButton{
			tgbotapi.NewInlineKeyboardButtonData(
				"🕵️ ناشناس",
				"comment_anonymous_"+
					strconv.FormatInt(
						teacherID,
						10,
					),
			),
		},

		[]tgbotapi.InlineKeyboardButton{
			tgbotapi.NewInlineKeyboardButtonData(
				"👤 با آیدی تلگرام",
				"comment_public_"+
					strconv.FormatInt(
						teacherID,
						10,
					),
			),
		},
	)
}