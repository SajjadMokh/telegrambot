package bot

import (
	"database/sql"
	"log"
	"strconv"
	"strings"

	"BOT/service"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type AddTeacherState struct {
	Step           int
	FirstName      string
	LastName       string
	Phone          string
	DepartmentName string
}

var teacherStates = make(map[int64]*AddTeacherState)

// =============================
// Handle Update
// =============================

func HandleUpdate(
	bot *tgbotapi.BotAPI,
	update tgbotapi.Update,
	db *sql.DB,
) {

	// =============================
	// Callback Query
	// =============================

	if update.CallbackQuery != nil {

		callback := update.CallbackQuery

		// پاسخ سریع به Telegram
		if _, err := bot.Request(
			tgbotapi.NewCallback(callback.ID, ""),
		); err != nil {
			log.Println("Callback answer error:", err)
		}

		if callback.Message == nil {
			return
		}

		// =============================
		// Teacher Profile
		// =============================

		if strings.HasPrefix(callback.Data, "teacher_") {

			teacherIDText := strings.TrimPrefix(
				callback.Data,
				"teacher_",
			)

			teacherID, err := strconv.ParseInt(
				teacherIDText,
				10,
				64,
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
				log.Println("Get teacher profile error:", err)

				msg := tgbotapi.NewMessage(
					callback.Message.Chat.ID,
					"❌ اطلاعات استاد پیدا نشد.",
				)

				_, _ = bot.Send(msg)

				return
			}

			// =============================
			// Rating
			// =============================

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

			// =============================
			// Teacher Profile Text
			// =============================

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

			msg := tgbotapi.NewMessage(
				callback.Message.Chat.ID,
				text,
			)

			// =============================
			// Profile Buttons
			// =============================

			keyboard := tgbotapi.NewInlineKeyboardMarkup(

				[]tgbotapi.InlineKeyboardButton{
					tgbotapi.NewInlineKeyboardButtonData(
						"⭐ ثبت امتیاز و نظر",
						"rate_"+stringID(teacherID),
					),
				},

				[]tgbotapi.InlineKeyboardButton{
					tgbotapi.NewInlineKeyboardButtonData(
						"💬 مشاهده نظرات",
						"comments_"+stringID(teacherID),
					),
				},
			)

			msg.ReplyMarkup = keyboard

			if _, err := bot.Send(msg); err != nil {
				log.Println(
					"Send teacher profile error:",
					err,
				)
			}

			return
		}

		// =============================
		// Rating
		// =============================

		if strings.HasPrefix(callback.Data, "rate_") {

			teacherIDText := strings.TrimPrefix(
				callback.Data,
				"rate_",
			)

			teacherID, err := strconv.ParseInt(
				teacherIDText,
				10,
				64,
			)

			if err != nil {
				log.Println("Invalid teacher ID:", err)
				return
			}

			msg := tgbotapi.NewMessage(
				callback.Message.Chat.ID,
				"⭐ امتیازت به این استاد رو از ۱ تا ۱۰ انتخاب کن:",
			)

			keyboard := ratingKeyboard(teacherID)

			msg.ReplyMarkup = keyboard

			if _, err := bot.Send(msg); err != nil {
				log.Println(
					"Send rating keyboard error:",
					err,
				)
			}

			return
		}

		// =============================
		// Comments
		// =============================

		if strings.HasPrefix(callback.Data, "comments_") {

			teacherIDText := strings.TrimPrefix(
				callback.Data,
				"comments_",
			)

			_, err := strconv.ParseInt(
				teacherIDText,
				10,
				64,
			)

			if err != nil {
				log.Println("Invalid teacher ID:", err)
				return
			}

			// فعلاً برای مرحله بعد
			msg := tgbotapi.NewMessage(
				callback.Message.Chat.ID,
				"💬 بخش نظرات رو در مرحله بعد کامل می‌کنیم.",
			)

			if _, err := bot.Send(msg); err != nil {
				log.Println(
					"Send comments message error:",
					err,
				)
			}

			return
		}

		// =============================
		// Rating Value
		// =============================

		if strings.HasPrefix(callback.Data, "rating_") {

			// فعلاً فقط دریافت Callback
			// منطق ثبت Rating را در مرحله بعد
			// به صورت کامل اضافه می‌کنیم.

			msg := tgbotapi.NewMessage(
				callback.Message.Chat.ID,
				"✅ امتیاز دریافت شد.\n\n"+
					"در مرحله بعد ثبت امتیاز و نظر را کامل می‌کنیم.",
			)

			if _, err := bot.Send(msg); err != nil {
				log.Println(
					"Send rating received message error:",
					err,
				)
			}

			return
		}

		return
	}

	// =============================
	// Ignore Non Message Updates
	// =============================

	if update.Message == nil {
		return
	}

	userID := update.Message.From.ID
	chatID := update.Message.Chat.ID

	text := strings.TrimSpace(
		update.Message.Text,
	)

	// =============================
	// /start
	// =============================

	if update.Message.IsCommand() &&
		update.Message.Command() == "start" {

		msg := tgbotapi.NewMessage(
			chatID,
			"سلام 👋\n\n"+
				"به ربات دانشگاه خوش اومدی 🌹\n\n"+
				"اسم استاد موردنظرت رو به فارسی بنویس:",
		)

		if _, err := bot.Send(msg); err != nil {
			log.Println(
				"Send start message error:",
				err,
			)
		}

		return
	}

	// =============================
	// /addteacher
	// =============================

	if update.Message.IsCommand() &&
		update.Message.Command() == "addteacher" {

		teacherService := service.NewTeacherService(db)

		if !teacherService.IsAdmin(userID) {

			msg := tgbotapi.NewMessage(
				chatID,
				"⛔ شما دسترسی افزودن استاد را ندارید.",
			)

			if _, err := bot.Send(msg); err != nil {
				log.Println(
					"Send admin error:",
					err,
				)
			}

			return
		}

		teacherStates[userID] = &AddTeacherState{
			Step: 1,
		}

		msg := tgbotapi.NewMessage(
			chatID,
			"➕ افزودن استاد\n\n"+
				"اسم استاد چیه؟",
		)

		if _, err := bot.Send(msg); err != nil {
			log.Println(
				"Send add teacher error:",
				err,
			)
		}

		return
	}

	// =============================
	// Add Teacher Process
	// =============================

	state, exists := teacherStates[userID]

	if exists {

		switch state.Step {

		// =============================
		// First Name
		// =============================

		case 1:

			if text == "" {
				return
			}

			state.FirstName = text
			state.Step = 2

			msg := tgbotapi.NewMessage(
				chatID,
				"نام خانوادگی استاد چیه؟",
			)

			_, _ = bot.Send(msg)

			return

		// =============================
		// Last Name
		// =============================

		case 2:

			if text == "" {
				return
			}

			state.LastName = text
			state.Step = 3

			msg := tgbotapi.NewMessage(
				chatID,
				"📞 شماره استاد چیه؟",
			)

			_, _ = bot.Send(msg)

			return

		// =============================
		// Phone
		// =============================

		case 3:

			if text == "" {
				return
			}

			state.Phone = text
			state.Step = 4

			msg := tgbotapi.NewMessage(
				chatID,
				"🏫 گروه یا دپارتمان استاد چیه؟\n\n"+
					"مثلاً:\n"+
					"مهندسی کامپیوتر",
			)

			_, _ = bot.Send(msg)

			return

		// =============================
		// Department
		// =============================

		case 4:

			if text == "" {
				return
			}

			state.DepartmentName = text

			teacherService := service.NewTeacherService(db)

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

				msg := tgbotapi.NewMessage(
					chatID,
					"❌ خطایی هنگام ثبت استاد اتفاق افتاد.",
				)

				_, _ = bot.Send(msg)

				delete(
					teacherStates,
					userID,
				)

				return
			}

			msg := tgbotapi.NewMessage(
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

			_, _ = bot.Send(msg)

			delete(
				teacherStates,
				userID,
			)

			return
		}
	}

	// =============================
	// Empty Text
	// =============================

	if text == "" {
		return
	}

	// =============================
	// Teacher Search
	// =============================

	searchService := service.NewSearchService(db)

	teachers, err := searchService.SearchTeachers(
		text,
	)

	if err != nil {

		log.Println(
			"Teacher search error:",
			err,
		)

		msg := tgbotapi.NewMessage(
			chatID,
			"❌ هنگام جستجوی استاد مشکلی پیش اومد.",
		)

		_, _ = bot.Send(msg)

		return
	}

	// =============================
	// No Result
	// =============================

	if len(teachers) == 0 {

		msg := tgbotapi.NewMessage(
			chatID,
			"😕 استادی با این نام پیدا نکردم.\n\n"+
				"لطفاً اسم یا نام خانوادگی استاد رو دوباره وارد کن.",
		)

		_, _ = bot.Send(msg)

		return
	}

	// =============================
	// Search Results
	// =============================

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

// =============================
// Rating Keyboard
// =============================

func ratingKeyboard(
	teacherID int64,
) tgbotapi.InlineKeyboardMarkup {

	var rows [][]tgbotapi.InlineKeyboardButton

	var row []tgbotapi.InlineKeyboardButton

	for rating := 1; rating <= 10; rating++ {

		button :=
			tgbotapi.NewInlineKeyboardButtonData(
				strconv.Itoa(rating)+" ⭐",
				"rating_"+
					stringID(teacherID)+
					"_"+
					strconv.Itoa(rating),
			)

		row = append(
			row,
			button,
		)

		// هر 5 امتیاز یک ردیف
		if len(row) == 5 {

			rows = append(
				rows,
				row,
			)

			row = nil
		}
	}

	if len(row) > 0 {
		rows = append(
			rows,
			row,
		)
	}

	return tgbotapi.NewInlineKeyboardMarkup(
		rows...,
	)
}

// =============================
// Convert ID to String
// =============================

func stringID(id int64) string {

	return strconv.FormatInt(
		id,
		10,
	)
}