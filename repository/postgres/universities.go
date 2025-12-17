package repository

import (
	"context"
	"database/sql"
	"db_labs/repository/postgres/stored"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type UniversitiesRepository struct {
	db *sqlx.DB
}

func NewUniversitiesRepository(db *sql.DB) *UniversitiesRepository {
	return &UniversitiesRepository{db: sqlx.NewDb(db, "pq")}
}

func (repo *UniversitiesRepository) Select(ctx context.Context) ([]stored.University, error) {
	query := "SELECT * from universities_select()"
	result := []stored.University{}
	err := repo.db.SelectContext(ctx, &result, query)
	return result, err
}

func (repo *UniversitiesRepository) Get(ctx context.Context, id uuid.UUID) (*stored.University, error) {
	var university = stored.University{Id: id}
	query := "universities_id_select($1)"
	err := repo.db.GetContext(ctx, university, query, id)
	return &university, err
}

func (repo *UniversitiesRepository) SelectByName(ctx context.Context, searched string) ([]stored.University, error) {
	query := "SELECT * from universities_name_fuzzy($1)"
	result := []stored.University{}
	err := repo.db.SelectContext(ctx, &result, query, searched)
	return result, err
}