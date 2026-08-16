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
	// Teacher Profile Callback
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

		// فقط Callbackهای مربوط به استاد
		if !strings.HasPrefix(callback.Data, "teacher_") {
			return
		}

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

		teacher, err := teacherService.GetTeacherByID(
			teacherID,
		)

		if err != nil {
			log.Println("Get teacher error:", err)

			msg := tgbotapi.NewMessage(
				callback.Message.Chat.ID,
				"❌ اطلاعات استاد پیدا نشد.",
			)

			bot.Send(msg)

			return
		}

		// =============================
		// Teacher Profile
		// =============================

		text := "👨‍🏫 پروفایل استاد\n\n" +
			"نام: " +
			teacher.FirstName +
			" " +
			teacher.LastName +
			"\n\n" +
			"📞 شماره تماس: " +
			teacher.Phone

		msg := tgbotapi.NewMessage(
			callback.Message.Chat.ID,
			text,
		)

		if _, err := bot.Send(msg); err != nil {
			log.Println("Send teacher profile error:", err)
		}

		return
	}

	if update.Message == nil {
		return
	}

	userID := update.Message.From.ID
	chatID := update.Message.Chat.ID
	text := strings.TrimSpace(update.Message.Text)

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
			log.Println("Send start message error:", err)
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
				log.Println("Send admin error:", err)
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
			log.Println("Send add teacher error:", err)
		}

		return
	}

	// =============================
	// Add Teacher Process
	// =============================

	state, exists := teacherStates[userID]

	if exists {

		switch state.Step {

		case 1:

			state.FirstName = text
			state.Step = 2

			msg := tgbotapi.NewMessage(
				chatID,
				"نام خانوادگی استاد چیه؟",
			)

			bot.Send(msg)

			return

		case 2:

			state.LastName = text
			state.Step = 3

			msg := tgbotapi.NewMessage(
				chatID,
				"📞 شماره استاد چیه؟",
			)

			bot.Send(msg)

			return

		case 3:

			state.Phone = text
			state.Step = 4

			msg := tgbotapi.NewMessage(
				chatID,
				"🏫 گروه یا دپارتمان استاد چیه؟\n\n"+
					"مثلاً:\n"+
					"مهندسی کامپیوتر",
			)

			bot.Send(msg)

			return

		case 4:

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

				log.Println("Add teacher error:", err)

				msg := tgbotapi.NewMessage(
					chatID,
					"❌ خطایی هنگام ثبت استاد اتفاق افتاد.",
				)

				bot.Send(msg)

				delete(teacherStates, userID)

				return
			}

			msg := tgbotapi.NewMessage(
				chatID,
				"✅ استاد با موفقیت ثبت شد.\n\n"+
					"👨‍🏫 "+state.FirstName+" "+state.LastName+"\n"+
					"📞 "+state.Phone+"\n"+
					"🏫 "+state.DepartmentName,
			)

			bot.Send(msg)

			delete(teacherStates, userID)

			return
		}
	}

	// =============================
	// Teacher Search
	// =============================

	if text == "" {
		return
	}

	searchService := service.NewSearchService(db)

	teachers, err := searchService.SearchTeachers(text)

	if err != nil {

		log.Println("Teacher search error:", err)

		msg := tgbotapi.NewMessage(
			chatID,
			"❌ هنگام جستجوی استاد مشکلی پیش اومد.",
		)

		bot.Send(msg)

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

		bot.Send(msg)

		return
	}

	// =============================
	// Search Results
	// =============================

	var rows [][]tgbotapi.InlineKeyboardButton

	for _, teacher := range teachers {

		fullName := teacher.FirstName + " " + teacher.LastName

		button := tgbotapi.NewInlineKeyboardButtonData(
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

	keyboard := tgbotapi.NewInlineKeyboardMarkup(rows...)

	msg := tgbotapi.NewMessage(
		chatID,
		"🔎 نتایج جستجو:\n\n"+
			"استاد موردنظرت رو انتخاب کن:",
	)

	msg.ReplyMarkup = keyboard

	if _, err := bot.Send(msg); err != nil {
		log.Println("Send search results error:", err)
	}
}

// =============================
// Convert ID to String
// =============================

func stringID(id int64) string {

	return strconv.FormatInt(id, 10)
}
