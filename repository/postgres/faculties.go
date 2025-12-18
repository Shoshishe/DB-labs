package postgres

import (
	"context"
	"database/sql"
	"db_labs/repository/postgres/stored"

	"github.com/jmoiron/sqlx"
)

type FacultiesRepository struct {
	db *sqlx.DB
}

func NewFacultiesRepository(db *sql.DB) *FacultiesRepository {
	return &FacultiesRepository{db: sqlx.NewDb(db, "pq")}
}

func (repo *FacultiesRepository) SelectFaculties(ctx context.Context) ([]stored.Faculty, error) {
	query := "SELECT * from faculties_select()"
	result := []stored.Faculty{}
	err := repo.db.SelectContext(ctx, &result, query)
	return result, err
}
