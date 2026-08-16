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
	Text        string
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
	// Comment Yes
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

		sendMessage(
			bot,
			chatID,
			"✍️ حالا نظرت درباره این استاد رو بنویس:",
		)

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

		state, exists := commentStates[userID]

		if !exists {
			log.Println(
				"Comment state not found for user:",
				userID,
			)
			return
		}

		// مطمئن می‌شویم کامنت مربوط به همان استاد است.
		if state.TeacherID != teacherID {
			log.Println(
				"Teacher ID mismatch in anonymous comment",
			)
			return
		}

		// باید متن نظر قبلاً دریافت شده باشد.
		if state.Step != 2 || state.Text == "" {
			log.Println("Comment text is not ready")
			return
		}

		state.IsAnonymous = true

		err = saveComment(
			bot,
			chatID,
			userID,
			state,
			db,
			callback.From.UserName,
		)

		if err != nil {
			log.Println(
				"Save anonymous comment error:",
				err,
			)
		}

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

		state, exists := commentStates[userID]

		if !exists {
			log.Println(
				"Comment state not found for user:",
				userID,
			)
			return
		}

		// مطمئن می‌شویم کامنت مربوط به همان استاد است.
		if state.TeacherID != teacherID {
			log.Println(
				"Teacher ID mismatch in public comment",
			)
			return
		}

		// باید متن نظر قبلاً دریافت شده باشد.
		if state.Step != 2 || state.Text == "" {
			log.Println("Comment text is not ready")
			return
		}

		state.IsAnonymous = false

		err = saveComment(
			bot,
			chatID,
			userID,
			state,
			db,
			callback.From.UserName,
		)

		if err != nil {
			log.Println(
				"Save public comment error:",
				err,
			)
		}

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

	// ========================================================
	// Step 1: Get Comment Text
	// ========================================================

	if state.Step != 1 {
		return false
	}

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
	// Save Text in State
	// ========================================================

	state.Text = text
	state.Step = 2

	// ========================================================
	// Ask Comment Identity
	// ========================================================

	msg := tgbotapi.NewMessage(
		message.Chat.ID,
		"💬 متن نظرت دریافت شد.\n\n"+
			"حالا انتخاب کن نظرت با چه حالتی ثبت بشه:",
	)

	msg.ReplyMarkup = commentIdentityKeyboard(
		state.TeacherID,
	)

	if _, err := bot.Send(msg); err != nil {

		log.Println(
			"Send comment identity keyboard error:",
			err,
		)
	}

	return true
}

// ============================================================
// Save Comment
// ============================================================

func saveComment(
	bot *tgbotapi.BotAPI,
	chatID int64,
	telegramUserID int64,
	state *CommentState,
	db *sql.DB,
	username string,
) error {

	// ========================================================
	// Get / Create User
	// ========================================================

	ratingService := service.NewRatingService(db)

	dbUserID, err := ratingService.GetOrCreateUser(
		telegramUserID,
		username,
	)

	if err != nil {

		log.Println(
			"Get or create user error:",
			err,
		)

		sendMessage(
			bot,
			chatID,
			"❌ خطایی هنگام شناسایی کاربر اتفاق افتاد.",
		)

		return err
	}

	// ========================================================
	// Save Comment
	// ========================================================

	commentService := service.NewCommentService(db)

	err = commentService.AddComment(
		state.TeacherID,
		dbUserID,
		state.Text,
		state.IsAnonymous,
	)

	if err != nil {

		log.Println(
			"Add comment error:",
			err,
		)

		sendMessage(
			bot,
			chatID,
			"❌ هنگام ثبت نظر مشکلی پیش اومد.",
		)

		return err
	}

	// ========================================================
	// Success
	// ========================================================

	sendMessage(
		bot,
		chatID,
		"✅ نظرت با موفقیت ثبت شد.\n\n"+
			"ممنون که تجربه‌ات رو به اشتراک گذاشتی ❤️",
	)

	// ========================================================
	// Clear State
	// ========================================================

	delete(
		commentStates,
		telegramUserID,
	)

	return nil
}

// ============================================================
// Comment Question Keyboard
// ============================================================

func commentQuestionKeyboard(
	teacherID int64,
) tgbotapi.InlineKeyboardMarkup {

	return tgbotapi.NewInlineKeyboardMarkup(

		[]tgbotapi.InlineKeyboardButton{
			tgbotapi.NewInlineKeyboardButtonData(
				"✅ بله، نظر می‌دم",
				"comment_yes_"+strconv.FormatInt(
					teacherID,
					10,
				),
			),
		},

		[]tgbotapi.InlineKeyboardButton{
			tgbotapi.NewInlineKeyboardButtonData(
				"❌ نه",
				"comment_no_"+strconv.FormatInt(
					teacherID,
					10,
				),
			),
		},
	)
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
				"comment_anonymous_"+strconv.FormatInt(
					teacherID,
					10,
				),
			),
		},

		[]tgbotapi.InlineKeyboardButton{
			tgbotapi.NewInlineKeyboardButtonData(
				"👤 با آیدی تلگرام",
				"comment_public_"+strconv.FormatInt(
					teacherID,
					10,
				),
			),
		},
	)
}