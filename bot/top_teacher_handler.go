package bot

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"BOT/service"
)

func HandleTopTeachersByMessage(
	bot *tgbotapi.BotAPI,
	chatID int64,
	db *sql.DB,
) {

	topService := service.NewTopTeacherService(db)

	teachers, err := topService.GetTopTeachers()

	if err != nil {

		log.Println(
			"Get top teachers error:",
			err,
		)

		sendMessage(
			bot,
			chatID,
			"❌ دریافت بهترین اساتید با مشکل مواجه شد.",
		)

		return
	}

	if len(teachers) == 0 {

		sendMessage(
			bot,
			chatID,
			"هنوز اطلاعات کافی برای رتبه‌بندی وجود ندارد.",
		)

		return
	}

	text := "🏆 رتبه‌بندی بهترین اساتید\n\n"

	for i, teacher := range teachers {

		text +=
			"🥇 رتبه " +
				strconv.Itoa(i+1) +
				"\n"

		text +=
			"👨‍🏫 " +
				teacher.FirstName +
				" " +
				teacher.LastName +
				"\n"

		text +=
			"🏫 " +
				teacher.DepartmentName +
				"\n"

		text +=
			"⭐ امتیاز نهایی: " +
				fmt.Sprintf(
					"%.2f",
					teacher.FinalScore,
				) +
				"\n"

		text +=
			"👥 تعداد رأی: " +
				strconv.Itoa(
					teacher.RatingCount,
				) +
				"\n\n────────────\n\n"

	}

	sendMessage(
		bot,
		chatID,
		text,
	)

}
