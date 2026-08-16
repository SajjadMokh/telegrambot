package service

import (
	"database/sql"
	"fmt"
	"strings"
	"unicode"

	"BOT/models"
)

type SearchService struct {
	DB *sql.DB
}

func NewSearchService(db *sql.DB) *SearchService {
	return &SearchService{
		DB: db,
	}
}

// =============================
// Normalize Persian Text
// =============================

func normalizePersianText(text string) string {

	text = strings.TrimSpace(text)

	// Arabic Yeh -> Persian Yeh
	text = strings.ReplaceAll(text, "ي", "ی")

	// Arabic Kaf -> Persian Kaf
	text = strings.ReplaceAll(text, "ك", "ک")

	// Remove Zero Width Non-Joiner
	text = strings.ReplaceAll(text, "\u200c", " ")

	// Normalize whitespace
	text = strings.Join(strings.Fields(text), " ")

	// Remove punctuation
	text = strings.Map(func(r rune) rune {

		if unicode.IsPunct(r) {
			return ' '
		}

		return r

	}, text)

	// Normalize whitespace again
	text = strings.Join(strings.Fields(text), " ")

	return text
}

// =============================
// Search Teachers
// =============================

func (s *SearchService) SearchTeachers(
	query string,
) ([]models.Teacher, error) {

	query = normalizePersianText(query)

	if query == "" {
		return []models.Teacher{}, nil
	}

	// حداقل دو کاراکتر
	if len([]rune(query)) < 2 {
		return []models.Teacher{}, nil
	}

	/*
		کاربر ممکن است بنویسد:

		علی احمدی
		علی    احمدی
		 علی احمدی
		علی احمدی!!!

		همه قبل از رسیدن به SQL نرمال شده‌اند.
	*/

	searchPattern := "%" + query + "%"
	prefixPattern := query + "%"

	rows, err := s.DB.Query(
		`
		SELECT
			t.id,
			t.first_name,
			t.last_name,
			t.phone,
			t.department_id,
			t.created_at

		FROM teachers t

		WHERE

			-- =========================
			-- Exact / Partial Full Name
			-- =========================

			CONCAT(t.first_name, ' ', t.last_name)
				ILIKE $1

			-- =========================
			-- First Name
			-- =========================

			OR t.first_name ILIKE $1

			-- =========================
			-- Last Name
			-- =========================

			OR t.last_name ILIKE $1

			-- =========================
			-- Fuzzy First Name
			-- =========================

			OR similarity(
				t.first_name,
				$2
			) >= 0.35

			-- =========================
			-- Fuzzy Last Name
			-- =========================

			OR similarity(
				t.last_name,
				$2
			) >= 0.35

		ORDER BY

			-- =========================
			-- 1. Exact Full Name
			-- =========================

			CASE
				WHEN CONCAT(
					t.first_name,
					' ',
					t.last_name
				) ILIKE $2
				THEN 1

			-- =========================
			-- 2. Full Name Starts With
			-- =========================

				WHEN CONCAT(
					t.first_name,
					' ',
					t.last_name
				) ILIKE $3
				THEN 2

			-- =========================
			-- 3. Exact Last Name
			-- =========================

				WHEN t.last_name ILIKE $2
				THEN 3

			-- =========================
			-- 4. Last Name Starts With
			-- =========================

				WHEN t.last_name ILIKE $3
				THEN 4

			-- =========================
			-- 5. First Name Starts With
			-- =========================

				WHEN t.first_name ILIKE $3
				THEN 5

			-- =========================
			-- 6. Last Name Contains
			-- =========================

				WHEN t.last_name ILIKE $1
				THEN 6

			-- =========================
			-- 7. First Name Contains
			-- =========================

				WHEN t.first_name ILIKE $1
				THEN 7

				ELSE 8
			END,

			-- =========================
			-- Fuzzy Similarity
			-- =========================

			GREATEST(
				similarity(t.first_name, $2),
				similarity(t.last_name, $2),
				similarity(
					CONCAT(
						t.first_name,
						' ',
						t.last_name
					),
					$2
				)
			) DESC,

			-- =========================
			-- Stable Ordering
			-- =========================

			t.last_name ASC,
			t.first_name ASC

		LIMIT 10
		`,
		searchPattern,
		query,
		prefixPattern,
	)

	if err != nil {
		return nil, fmt.Errorf(
			"search teachers: %w",
			err,
		)
	}

	defer rows.Close()

	teachers := make(
		[]models.Teacher,
		0,
	)

	for rows.Next() {

		var teacher models.Teacher

		err := rows.Scan(
			&teacher.ID,
			&teacher.FirstName,
			&teacher.LastName,
			&teacher.Phone,
			&teacher.DepartmentID,
			&teacher.CreatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf(
				"scan teacher: %w",
				err,
			)
		}

		teachers = append(
			teachers,
			teacher,
		)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf(
			"iterate teachers: %w",
			err,
		)
	}

	return teachers, nil
}