package games

var CurrentQuestion = map[int64]int{}

var Scores = map[int64]int{}


// شروع بازی برای کاربر
func StartQuiz(userID int64) {

	CurrentQuestion[userID] = 0
	Scores[userID] = 0

}


// گرفتن شماره سوال فعلی
func GetQuestionIndex(userID int64) int {

	return CurrentQuestion[userID]

}


// رفتن به سوال بعدی
func NextQuestion(userID int64) {

	CurrentQuestion[userID]++

}


// افزایش امتیاز
func AddScore(userID int64) {

	Scores[userID]++

}


// گرفتن امتیاز
func GetScore(userID int64) int {

	return Scores[userID]

}