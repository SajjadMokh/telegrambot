package service

import (
	"database/sql"
	"fmt"
	"strings"
)

type CommentService struct {
	DB *sql.DB
}

func NewCommentService(db *sql.DB) *CommentService {
	return &CommentService{
		DB: db,
	}
}

// =============================
// Add Comment
// =============================

func (s *CommentService) AddComment(
	teacherID int64,
	userID int64,
	text string,
	isAnonymous bool,
) error {

	text = strings.TrimSpace(text)

	if text == "" {
		return fmt.Errorf("comment text is empty")
	}

	_, err := s.DB.Exec(
		`
        INSERT INTO comments (
            teacher_id,
            user_id,
            text,
            is_anonymous
        )
        VALUES ($1, $2, $3, $4)
        `,
		teacherID,
		userID,
		text,
		isAnonymous,
	)

	if err != nil {
		return fmt.Errorf("add comment: %w", err)
	}

	return nil
}
