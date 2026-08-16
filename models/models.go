package models

import "time"

// =============================
// Department
// =============================

type Department struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

// =============================
// Teacher
// =============================

type Teacher struct {
	ID           int64     `json:"id"`
	FirstName    string    `json:"first_name"`
	LastName     string    `json:"last_name"`
	Phone        string    `json:"phone"`
	DepartmentID int64     `json:"department_id"`
	CreatedAt    time.Time `json:"created_at"`
}
// =============================
// User
// =============================

type User struct {
	ID         int64     `json:"id"`
	TelegramID int64     `json:"telegram_id"`
	Username   string    `json:"username"`
	CreatedAt  time.Time `json:"created_at"`
}

// =============================
// Rating
// =============================

type Rating struct {
	ID        int64     `json:"id"`
	TeacherID int64     `json:"teacher_id"`
	UserID    int64     `json:"user_id"`
	Rating    int       `json:"rating"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// =============================
// Comment
// =============================

type Comment struct {
	ID          int64     `json:"id"`
	TeacherID   int64     `json:"teacher_id"`
	UserID      int64     `json:"user_id"`
	Text        string    `json:"text"`
	IsAnonymous bool      `json:"is_anonymous"`
	CreatedAt   time.Time `json:"created_at"`
}