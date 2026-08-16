package service

import (
	"BOT/models"
	"database/sql"
	"fmt"
	"strings"
)

const AdminTelegramID int64 = 1150702474

type TeacherService struct {
	DB *sql.DB
}

func NewTeacherService(db *sql.DB) *TeacherService {
	return &TeacherService{
		DB: db,
	}
}

// =============================
// Check Admin
// =============================

func (s *TeacherService) IsAdmin(telegramID int64) bool {
	return telegramID == AdminTelegramID
}

// =============================
// Add Teacher
// =============================

func (s *TeacherService) AddTeacher(
	telegramID int64,
	firstName string,
	lastName string,
	phone string,
	departmentName string,
) error {

	if !s.IsAdmin(telegramID) {
		return fmt.Errorf("access denied")
	}

	firstName = strings.TrimSpace(firstName)
	lastName = strings.TrimSpace(lastName)
	phone = strings.TrimSpace(phone)
	departmentName = strings.TrimSpace(departmentName)

	if firstName == "" {
		return fmt.Errorf("first name is required")
	}

	if lastName == "" {
		return fmt.Errorf("last name is required")
	}

	if phone == "" {
		return fmt.Errorf("phone is required")
	}

	if departmentName == "" {
		return fmt.Errorf("department is required")
	}

	tx, err := s.DB.Begin()
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}

	defer tx.Rollback()

	// =============================
	// Find or Create Department
	// =============================

	var departmentID int64

	err = tx.QueryRow(
		`
		INSERT INTO departments (name)
		VALUES ($1)
		ON CONFLICT (name)
		DO UPDATE SET name = EXCLUDED.name
		RETURNING id
		`,
		departmentName,
	).Scan(&departmentID)

	if err != nil {
		return fmt.Errorf("get department: %w", err)
	}

	// =============================
	// Insert Teacher
	// =============================

	_, err = tx.Exec(
		`
		INSERT INTO teachers
		(
			first_name,
			last_name,
			phone,
			department_id
		)
		VALUES ($1, $2, $3, $4)
		`,
		firstName,
		lastName,
		phone,
		departmentID,
	)

	if err != nil {
		return fmt.Errorf("insert teacher: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}

// =============================
// Get Teacher By ID
// =============================

func (s *TeacherService) GetTeacherByID(
	teacherID int64,
) (*models.Teacher, error) {

	var teacher models.Teacher

	err := s.DB.QueryRow(
		`
		SELECT
			t.id,
			t.first_name,
			t.last_name,
			t.phone,
			t.department_id,
			t.created_at
		FROM teachers t
		WHERE t.id = $1
		`,
		teacherID,
	).Scan(
		&teacher.ID,
		&teacher.FirstName,
		&teacher.LastName,
		&teacher.Phone,
		&teacher.DepartmentID,
		&teacher.CreatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf(
			"get teacher by id: %w",
			err,
		)
	}

	return &teacher, nil
}

// =============================
// Get Teacher Profile
// =============================

func (s *TeacherService) GetTeacherProfile(
	teacherID int64,
) (*models.TeacherProfile, error) {

	var profile models.TeacherProfile

	err := s.DB.QueryRow(
		`
		SELECT
			t.id,
			t.first_name,
			t.last_name,
			t.phone,
			t.department_id,
			d.name AS department_name,

			COALESCE(AVG(r.rating), 0)::FLOAT8
				AS average_rating,

			COUNT(r.id)::INT
				AS rating_count

		FROM teachers t

		INNER JOIN departments d
			ON d.id = t.department_id

		LEFT JOIN ratings r
			ON r.teacher_id = t.id

		WHERE t.id = $1

		GROUP BY
			t.id,
			t.first_name,
			t.last_name,
			t.phone,
			t.department_id,
			d.name
		`,
		teacherID,
	).Scan(
		&profile.ID,
		&profile.FirstName,
		&profile.LastName,
		&profile.Phone,
		&profile.DepartmentID,
		&profile.DepartmentName,
		&profile.AverageRating,
		&profile.RatingCount,
	)

	if err != nil {
		return nil, fmt.Errorf(
			"get teacher profile: %w",
			err,
		)
	}

	return &profile, nil
}

// =============================
// Get Department Name
// =============================

func (s *TeacherService) GetDepartmentName(
	departmentID int64,
) (string, error) {

	var name string

	err := s.DB.QueryRow(
		`
		SELECT name
		FROM departments
		WHERE id = $1
		`,
		departmentID,
	).Scan(&name)

	if err != nil {
		return "", fmt.Errorf(
			"get department name: %w",
			err,
		)
	}

	return name, nil
}