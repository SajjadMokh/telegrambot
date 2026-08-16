package service

import (
	"database/sql"
	"fmt"
)

// ============================================================
// Comment
// ============================================================

type Comment struct {
	ID          int64
	TeacherID   int64
	UserID      int64
	Username    string
	Text        string
	IsAnonymous bool
}

type CommentService struct {
	DB *sql.DB
}

func NewCommentService(db *sql.DB) *CommentService {
	return &CommentService{
		DB: db,
	}
}

// ============================================================
// Add Comment
// ============================================================

func (s *CommentService) AddComment(
	teacherID int64,
	userID int64,
	text string,
	isAnonymous bool,
) error {

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
		return fmt.Errorf(
			"add comment: %w",
			err,
		)
	}

	return nil
}

// ============================================================
// Get Teacher Comments
// ============================================================

func (s *CommentService) GetTeacherComments(
	teacherID int64,
) ([]Comment, error) {

	rows, err := s.DB.Query(
		`
		SELECT
			c.id,
			c.teacher_id,
			c.user_id,
			COALESCE(u.username, ''),
			c.text,
			c.is_anonymous
		FROM comments c
		LEFT JOIN users u
			ON u.id = c.user_id
		WHERE c.teacher_id = $1
		ORDER BY c.id DESC
		`,
		teacherID,
	)

	if err != nil {
		return nil, fmt.Errorf(
			"get teacher comments: %w",
			err,
		)
	}

	defer rows.Close()

	var comments []Comment

	for rows.Next() {

		var comment Comment

		err := rows.Scan(
			&comment.ID,
			&comment.TeacherID,
			&comment.UserID,
			&comment.Username,
			&comment.Text,
			&comment.IsAnonymous,
		)

		if err != nil {
			return nil, fmt.Errorf(
				"scan comment: %w",
				err,
			)
		}

		comments = append(
			comments,
			comment,
		)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf(
			"iterate comments: %w",
			err,
		)
	}

	return comments, nil
}