package postgres

import (
	"context"
	"database/sql"
	"db_labs/repository/postgres/stored"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type LessonsRepository struct {
	db *sqlx.DB
}

func NewLessonsRepository(db *sql.DB) *LessonsRepository {
	return &LessonsRepository{db: sqlx.NewDb(db, "pq")}
}

func (repo *LessonsRepository) GroupGradesOrdered(ctx context.Context, groupId uuid.UUID) ([]stored.GroupGrades, error) {
	query := "SELECT * FROM group_grades_ordered($1)"
	result := []stored.GroupGrades{}
	err := repo.db.SelectContext(ctx, result, query)
	return result, err
}

func (repo *LessonsRepository) SelectOverGrade(ctx context.Context, grade float64) ([]stored.GroupGrades, error) {
	query := "SELECT * FROM grades_from($1)"
	result := []stored.GroupGrades{}
	err := repo.db.SelectContext(ctx, result, query)
	return result, err
}
