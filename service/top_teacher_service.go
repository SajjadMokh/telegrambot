package service

import (
	"database/sql"
	"fmt"

	"BOT/models"
)

type TopTeacherService struct {
	db *sql.DB
}

func NewTopTeacherService(
	db *sql.DB,
) *TopTeacherService {

	return &TopTeacherService{
		db: db,
	}
}

// ============================================================
// Get Top Teachers
// ============================================================

func (s *TopTeacherService) GetTopTeachers() ([]models.Teacher, error) {

	query := `
WITH global_average AS (

	SELECT 
		COALESCE(AVG(rating),0) AS c
	FROM ratings

),

teacher_scores AS (

	SELECT

		t.id,
		t.first_name,
		t.last_name,
		t.phone,
		t.department_id,
		d.name AS department_name,
		t.created_at,


		COALESCE(AVG(r.rating),0) AS average_rating,


		COUNT(r.id) AS rating_count,


		(
			(
				COUNT(r.id)::float 
				/
				(COUNT(r.id) + 5)
			)
			*
			COALESCE(AVG(r.rating),0)

		)

		+

		(
			(
				5::float
				/
				(COUNT(r.id)+5)
			)

			*
			global_average.c
		)

		AS final_score


	FROM teachers t


	LEFT JOIN departments d
	ON d.id = t.department_id


	LEFT JOIN ratings r
	ON r.teacher_id = t.id


	CROSS JOIN global_average


	GROUP BY

		t.id,
		d.name,
		global_average.c

)


SELECT

	id,
	first_name,
	last_name,
	phone,
	department_id,
	department_name,
	average_rating,
	rating_count,
	created_at,
	final_score


FROM teacher_scores


WHERE rating_count > 0


ORDER BY final_score DESC


LIMIT 3;

`

	rows, err := s.db.Query(query)

	if err != nil {
		return nil, fmt.Errorf("top teacher query failed: %w", err)
	}
	defer rows.Close()

	var teachers []models.Teacher

	for rows.Next() {

		var teacher models.Teacher

		err := rows.Scan(

			&teacher.ID,

			&teacher.FirstName,

			&teacher.LastName,

			&teacher.Phone,

			&teacher.DepartmentID,

			&teacher.DepartmentName,

			&teacher.AverageRating,

			&teacher.RatingCount,

			&teacher.CreatedAt,

			&teacher.FinalScore,
		)

		if err != nil {
			return nil, err
		}

		teachers = append(
			teachers,
			teacher,
		)

		if err := rows.Err(); err != nil {
			return nil, fmt.Errorf("top teacher rows error: %w", err)
		}
	}

	return teachers, nil
}
