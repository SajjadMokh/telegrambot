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
// Add Teacher State
// ============================================================

type AddTeacherState struct {
	Step           int
	FirstName      string
	LastName       string
	Phone          string
	DepartmentName string
}

var teacherStates = make(map[int64]*AddTeacherState)

// ============================================================
// Handle Update
// ============================================================

func HandleUpdate(
	bot *tgbotapi.BotAPI,
	update tgbotapi.Update,
	db *sql.DB,
) {

	// ========================================================
	// Callback Query
	// ========================================================

	if update.CallbackQuery != nil {

		callback := update.CallbackQuery

		// پاسخ سریع به Callback تلگرام
		if _, err := bot.Request(
			tgbotapi.NewCallback(callback.ID, ""),
		); err != nil {
			log.Println("Callback answer error:", err)
		}

		if callback.Message == nil {
			return
		}

		chatID := callback.Message.Chat.ID
		data := callback.Data

		// ====================================================
		// Teacher Profile
		// ====================================================

		if strings.HasPrefix(data, "teacher_") {

			teacherID, err := parseID(
				data,
				"teacher_",
			)

			if err != nil {
				log.Println("Invalid teacher ID:", err)
				return
			}

			teacherService := service.NewTeacherService(db)

			profile, err := teacherService.GetTeacherProfile(
				teacherID,
			)

			if err != nil {

				log.Println(
					"Get teacher profile error:",
					err,
				)

				sendMessage(
					bot,
					chatID,
					"❌ اطلاعات استاد پیدا نشد.",
				)

				return
			}

			// =================================================
			// Rating Text
			// =================================================

			ratingText := "هنوز امتیازی ثبت نشده ⭐"

			if profile.RatingCount > 0 {

				ratingText =
					strconv.FormatFloat(
						profile.AverageRating,
						'f',
						1,
						64,
					) + " / 10 ⭐"
			}

			// =================================================
			// Teacher Profile Text
			// =================================================

			text :=
				"👨‍🏫 پروفایل استاد\n\n" +
					"👤 نام: " +
					profile.FirstName +
					" " +
					profile.LastName +
					"\n\n" +
					"🏫 گروه: " +
					profile.DepartmentName +
					"\n\n" +
					"📞 شماره تماس: " +
					profile.Phone +
					"\n\n" +
					"⭐ امتیاز: " +
					ratingText +
					"\n" +
					"👥 تعداد رأی‌ها: " +
					strconv.Itoa(profile.RatingCount)

			// =================================================
			// Teacher Profile Keyboard
			// =================================================

			msg := tgbotapi.NewMessage(
				chatID,
				text,
			)

			msg.ReplyMarkup = teacherProfileKeyboard(
				teacherID,
			)

			if _, err := bot.Send(msg); err != nil {

				log.Println(
					"Send teacher profile error:",
					err,
				)
			}

			return
		}

		// ====================================================
		// Rating Callback
		// ====================================================

		if strings.HasPrefix(data, "rate_") ||
			strings.HasPrefix(data, "rating_") {

			HandleRatingCallback(
				bot,
				callback,
				db,
			)

			return
		}

		// ====================================================
		// Show Comments
		// ====================================================

		if strings.HasPrefix(data, "comments_") {

			teacherID, err := parseID(
				data,
				"comments_",
			)

			if err != nil {
				log.Println("Invalid teacher ID:", err)
				return
			}

			commentService := service.NewCommentService(db)

			comments, err := commentService.GetTeacherComments(
				teacherID,
			)

			if err != nil {

				log.Println(
					"Get teacher comments error:",
					err,
				)

				sendMessage(
					bot,
					chatID,
					"❌ دریافت نظرات با خطا مواجه شد.",
				)

				return
			}

			// =================================================
			// No Comments
			// =================================================

			if len(comments) == 0 {

				sendMessage(
					bot,
					chatID,
					"💬 هنوز نظری برای این استاد ثبت نشده.",
				)

				return
			}

			// =================================================
			// Build Comments
			// =================================================

			var builder strings.Builder

			builder.WriteString(
				"💬 نظرات درباره استاد\n\n",
			)

			for i, comment := range comments {

				builder.WriteString(
					"💬 نظر ",
				)

				builder.WriteString(
					strconv.Itoa(i + 1),
				)

				builder.WriteString("\n")

				// =============================================
				// Anonymous / Public
				// =============================================

				if comment.IsAnonymous {

					// اگر ناشناس باشد، Username اصلاً نمایش داده نمی‌شود.
					builder.WriteString(
						"👤 ناشناس\n",
					)

				} else {

					username := strings.TrimSpace(
						comment.Username,
					)

					if username != "" {

						builder.WriteString(
							"👤 @",
						)

						builder.WriteString(
							username,
						)

						builder.WriteString("\n")

					} else {

						builder.WriteString(
							"👤 کاربر تلگرام\n",
						)
					}
				}

				// =============================================
				// Comment Text
				// =============================================

				builder.WriteString(
					comment.Text,
				)

				builder.WriteString(
					"\n\n────────────\n\n",
				)
			}

			sendMessage(
				bot,
				chatID,
				builder.String(),
			)

			return
		}

		// ====================================================
		// Comment Callback
		// ====================================================

		if strings.HasPrefix(data, "comment_") {

			HandleCommentCallback(
				bot,
				callback,
				db,
			)

			return
		}

		// ====================================================
		// Unknown Callback
		// ====================================================

		log.Println(
			"Unknown callback:",
			data,
		)

		return
	}

	// ========================================================
	// Ignore Non Message Updates
	// ========================================================

	if update.Message == nil {
		return
	}

	message := update.Message

	userID := message.From.ID
	chatID := message.Chat.ID
	text := strings.TrimSpace(message.Text)
	// ========================================================
	// Top Teachers Menu
	// ========================================================

	if text == "🏆 بهترین استاد" {

		HandleTopTeachersByMessage(
			bot,
			chatID,
			db,
		)

		return
	}
	// ========================================================
	// Main Menu Actions
	// ========================================================

	if text == "🏆 بهترین استاد" {

		msg := tgbotapi.NewMessage(
			chatID,
			"🏆 سه استاد برتر دانشگاه:\n\n",
		)

		if _, err := bot.Send(msg); err != nil {
			log.Println(
				"Send top teacher request error:",
				err,
			)
		}

		HandleTopTeachersByMessage(
			bot,
			chatID,
			db,
		)

		return
	}

	if text == "🔎 جستجوی استاد" {

		sendMessage(
			bot,
			chatID,
			"🔎 اسم یا نام خانوادگی استاد را وارد کنید:",
		)

		return
	}

	// ========================================================
	// Comment Message
	// ========================================================

	if HandleCommentMessage(
		bot,
		message,
		db,
	) {
		return
	}

	// ========================================================
	// /start
	// ========================================================

	if message.IsCommand() &&
		message.Command() == "start" {

		msg := tgbotapi.NewMessage(
			chatID,
			"سلام دانشجوی عزیز 👋\n\n"+
				"به ربات ارزیابی اساتید خوش آمدید 🌹\n\n"+
				"یکی از گزینه‌های زیر را انتخاب کنید:",
		)

		msg.ReplyMarkup = mainMenuKeyboard()

		if _, err := bot.Send(msg); err != nil {
			log.Println(
				"Send start menu error:",
				err,
			)
		}

		return
	}

	// ========================================================
	// /addteacher
	// ========================================================

	if message.IsCommand() &&
		message.Command() == "addteacher" {

		teacherService :=
			service.NewTeacherService(db)

		if !teacherService.IsAdmin(userID) {

			sendMessage(
				bot,
				chatID,
				"⛔ شما دسترسی افزودن استاد را ندارید.",
			)

			return
		}

		teacherStates[userID] =
			&AddTeacherState{
				Step: 1,
			}

		sendMessage(
			bot,
			chatID,
			"➕ افزودن استاد\n\n"+
				"اسم استاد چیه؟",
		)

		return
	}

	// ========================================================
	// Add Teacher Process
	// ========================================================

	if state, exists := teacherStates[userID]; exists {

		switch state.Step {

		// ====================================================
		// First Name
		// ====================================================

		case 1:

			if text == "" {
				return
			}

			state.FirstName = text
			state.Step = 2

			sendMessage(
				bot,
				chatID,
				"نام خانوادگی استاد چیه؟",
			)

			return

		// ====================================================
		// Last Name
		// ====================================================

		case 2:

			if text == "" {
				return
			}

			state.LastName = text
			state.Step = 3

			sendMessage(
				bot,
				chatID,
				"📞 شماره استاد چیه؟",
			)

			return

		// ====================================================
		// Phone
		// ====================================================

		case 3:

			if text == "" {
				return
			}

			state.Phone = text
			state.Step = 4

			sendMessage(
				bot,
				chatID,
				"🏫 گروه یا دپارتمان استاد چیه؟\n\n"+
					"مثلاً:\n"+
					"مهندسی کامپیوتر",
			)

			return

		// ====================================================
		// Department
		// ====================================================

		case 4:

			if text == "" {
				return
			}

			state.DepartmentName = text

			teacherService :=
				service.NewTeacherService(db)

			err := teacherService.AddTeacher(
				userID,
				state.FirstName,
				state.LastName,
				state.Phone,
				state.DepartmentName,
			)

			if err != nil {

				log.Println(
					"Add teacher error:",
					err,
				)

				sendMessage(
					bot,
					chatID,
					"❌ خطایی هنگام ثبت استاد اتفاق افتاد.",
				)

				delete(
					teacherStates,
					userID,
				)

				return
			}

			sendMessage(
				bot,
				chatID,
				"✅ استاد با موفقیت ثبت شد.\n\n"+
					"👨‍🏫 "+
					state.FirstName+
					" "+
					state.LastName+
					"\n"+
					"📞 "+
					state.Phone+
					"\n"+
					"🏫 "+
					state.DepartmentName,
			)

			delete(
				teacherStates,
				userID,
			)

			return
		}
	}

	// ========================================================
	// Empty Text
	// ========================================================

	if text == "" {
		return
	}

	// ========================================================
	// Teacher Search
	// ========================================================

	searchService :=
		service.NewSearchService(db)

	teachers, err :=
		searchService.SearchTeachers(text)

	if err != nil {

		log.Println(
			"Teacher search error:",
			err,
		)

		sendMessage(
			bot,
			chatID,
			"❌ هنگام جستجوی استاد مشکلی پیش اومد.",
		)

		return
	}

	// ========================================================
	// No Result
	// ========================================================

	if len(teachers) == 0 {

		sendMessage(
			bot,
			chatID,
			"😕 استادی با این نام پیدا نکردم.\n\n"+
				"لطفاً اسم یا نام خانوادگی استاد رو دوباره وارد کن.",
		)

		return
	}

	// ========================================================
	// Search Results
	// ========================================================

	var rows [][]tgbotapi.InlineKeyboardButton

	for _, teacher := range teachers {

		fullName :=
			teacher.FirstName +
				" " +
				teacher.LastName

		button :=
			tgbotapi.NewInlineKeyboardButtonData(
				"👨‍🏫 "+fullName,
				"teacher_"+stringID(teacher.ID),
			)

		rows = append(
			rows,
			[]tgbotapi.InlineKeyboardButton{
				button,
			},
		)
	}

	keyboard :=
		tgbotapi.NewInlineKeyboardMarkup(
			rows...,
		)

	msg := tgbotapi.NewMessage(
		chatID,
		"🔎 نتایج جستجو:\n\n"+
			"استاد موردنظرت رو انتخاب کن:",
	)

	msg.ReplyMarkup = keyboard

	if _, err := bot.Send(msg); err != nil {

		log.Println(
			"Send search results error:",
			err,
		)
	}
}

// ============================================================
// Parse ID
// ============================================================

func parseID(
	data string,
	prefix string,
) (int64, error) {

	value := strings.TrimPrefix(
		data,
		prefix,
	)

	return strconv.ParseInt(
		value,
		10,
		64,
	)
}

// ============================================================
// Send Message
// ============================================================

func sendMessage(
	bot *tgbotapi.BotAPI,
	chatID int64,
	text string,
) {

	msg := tgbotapi.NewMessage(
		chatID,
		text,
	)

	if _, err := bot.Send(msg); err != nil {

		log.Println(
			"Send message error:",
			err,
		)
	}
}

// ============================================================
// Convert ID to String
// ============================================================

func stringID(id int64) string {

	return strconv.FormatInt(
		id,
		10,
	)
}
