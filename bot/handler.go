package bot

import (
	"database/sql"
	"log"
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

	if update.Message == nil {
		return
	}

	userID := update.Message.From.ID
	text := strings.TrimSpace(update.Message.Text)

	// =============================
	// /start
	// =============================

	if update.Message.IsCommand() &&
		update.Message.Command() == "start" {

		msg := tgbotapi.NewMessage(
			update.Message.Chat.ID,
			"سلام 👋\n\n"+
				"به ربات دانشگاه خوش اومدی 🌹\n\n"+
				"🔎 اسم استاد موردنظرت رو بفرست تا برات پیداش کنم.",
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
				update.Message.Chat.ID,
				"⛔ شما دسترسی افزودن استاد را ندارید.",
			)

			bot.Send(msg)

			return
		}

		teacherStates[userID] = &AddTeacherState{
			Step: 1,
		}

		msg := tgbotapi.NewMessage(
			update.Message.Chat.ID,
			"➕ افزودن استاد\n\n"+
				"اسم استاد چیه؟",
		)

		bot.Send(msg)

		return
	}

	// =============================
	// Add Teacher Process
	// =============================

	state, exists := teacherStates[userID]

	if !exists {
		return
	}

	switch state.Step {

	case 1:

		state.FirstName = text
		state.Step = 2

		msg := tgbotapi.NewMessage(
			update.Message.Chat.ID,
			"نام خانوادگی استاد چیه؟",
		)

		bot.Send(msg)

	case 2:

		state.LastName = text
		state.Step = 3

		msg := tgbotapi.NewMessage(
			update.Message.Chat.ID,
			"📞 شماره استاد چیه؟",
		)

		bot.Send(msg)

	case 3:

		state.Phone = text
		state.Step = 4

		msg := tgbotapi.NewMessage(
			update.Message.Chat.ID,
			"🏫 گروه یا دپارتمان استاد چیه؟\n\n"+
				"مثلاً:\n"+
				"مهندسی کامپیوتر",
		)

		bot.Send(msg)

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
				update.Message.Chat.ID,
				"❌ خطایی هنگام ثبت استاد اتفاق افتاد.\n\n"+
					"اطلاعات را بررسی کن و دوباره تلاش کن.",
			)

			bot.Send(msg)

			delete(teacherStates, userID)

			return
		}

		msg := tgbotapi.NewMessage(
			update.Message.Chat.ID,
			"✅ استاد با موفقیت ثبت شد.\n\n"+
				"👨‍🏫 "+state.FirstName+" "+state.LastName+"\n"+
				"📞 "+state.Phone+"\n"+
				"🏫 "+state.DepartmentName,
		)

		bot.Send(msg)

		delete(teacherStates, userID)
	}
}