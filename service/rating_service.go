package service

import (
	"database/sql"
	"fmt"
)

type RatingService struct {
	DB *sql.DB
}

func NewRatingService(db *sql.DB) *RatingService {
	return &RatingService{
		DB: db,
	}
}

// =============================
// Get or Create User
// =============================

func (s *RatingService) GetOrCreateUser(
	telegramID int64,
	username string,
) (int64, error) {

	var userID int64

	err := s.DB.QueryRow(
		`
		INSERT INTO users (
			telegram_id,
			username
		)
		VALUES ($1, $2)
		ON CONFLICT (telegram_id)
		DO UPDATE SET username = EXCLUDED.username
		RETURNING id
		`,
		telegramID,
		username,
	).Scan(&userID)

	if err != nil {
		return 0, fmt.Errorf(
			"get or create user: %w",
			err,
		)
	}

	return userID, nil
}

// =============================
// Save Rating
// =============================

func (s *RatingService) SaveRating(
	userID int64,
	teacherID int64,
	rating int,
) error {

	if rating < 1 || rating > 10 {
		return fmt.Errorf("rating must be between 1 and 10")
	}

	_, err := s.DB.Exec(
		`
		INSERT INTO ratings (
			teacher_id,
			user_id,
			rating
		)
		VALUES ($1, $2, $3)

		ON CONFLICT (user_id, teacher_id)
		DO UPDATE SET
			rating = EXCLUDED.rating,
			updated_at = NOW()
		`,
		teacherID,
		userID,
		rating,
	)

	if err != nil {
		return fmt.Errorf(
			"save rating: %w",
			err,
		)
	}

	return nil
}

// =============================
// Get Teacher Rating
// =============================

func (s *RatingService) GetTeacherRating(
	teacherID int64,
) (float64, int, error) {

	var average float64
	var count int

	err := s.DB.QueryRow(
		`
		SELECT
			COALESCE(AVG(rating), 0),
			COUNT(*)
		FROM ratings
		WHERE teacher_id = $1
		`,
		teacherID,
	).Scan(
		&average,
		&count,
	)

	if err != nil {
		return 0, 0, fmt.Errorf(
			"get teacher rating: %w",
			err,
		)
	}

	return average, count, nil
}