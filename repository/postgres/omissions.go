package postgres

import (
	"context"
	"database/sql"
	"db_labs/repository/postgres/stored"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type OmmissionsRepository struct {
	db *sqlx.DB
}

func NewOmissionsRepository(db *sql.DB) *OmmissionsRepository {
	return &OmmissionsRepository{db: sqlx.NewDb(db, "pq")}
}

func (repo *OmmissionsRepository) GetStudentOmissions(ctx context.Context, studentId uuid.UUID) ([]stored.Omission, error) {
	query := "SELECT * from get_student_omissions($1)"
	result := []stored.Omission{}
	err := repo.db.SelectContext(ctx, &result, query, studentId)
	return result, err
}
