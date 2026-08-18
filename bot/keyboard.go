package bot

import (
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// ============================================================
// Main Menu Keyboard
// ============================================================

func mainMenuKeyboard() tgbotapi.InlineKeyboardMarkup {

	return tgbotapi.NewInlineKeyboardMarkup(

		[]tgbotapi.InlineKeyboardButton{
			tgbotapi.NewInlineKeyboardButtonData(
				"🔎 جستجوی استاد",
				"search_teacher",
			),
		},

		[]tgbotapi.InlineKeyboardButton{
			tgbotapi.NewInlineKeyboardButtonData(
				"🏆 برترین استادها",
				"top_teachers",
			),
		},
	)
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
				"_"+
				strconv.Itoa(rating),
		)

		row = append(row, button)

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


func teacherProfileKeyboard(
	teacherID int64,
) tgbotapi.InlineKeyboardMarkup {

	return tgbotapi.NewInlineKeyboardMarkup(

		[]tgbotapi.InlineKeyboardButton{
			tgbotapi.NewInlineKeyboardButtonData(
				"⭐ ثبت امتیاز",
				"rate_"+strconv.FormatInt(
					teacherID,
					10,
				),
			),
		},

		[]tgbotapi.InlineKeyboardButton{
			tgbotapi.NewInlineKeyboardButtonData(
				"💬 ثبت نظر",
				"comment_start_"+strconv.FormatInt(
					teacherID,
					10,
				),
			),
		},

		[]tgbotapi.InlineKeyboardButton{
			tgbotapi.NewInlineKeyboardButtonData(
				"📖 مشاهده نظرات",
				"comments_"+strconv.FormatInt(
					teacherID,
					10,
				),
			),
		},
	)
}