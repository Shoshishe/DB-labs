package postgres

import (
	"context"
	"database/sql"
	"db_labs/repository/postgres/stored"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type GroupsRepository struct {
	db *sqlx.DB
}

func NewGroupsRepository(db *sql.DB) *GroupsRepository {
	return &GroupsRepository{db: sqlx.NewDb(db, "pq")}
}

func (repo *GroupsRepository) SelectGroupLessonNames(ctx context.Context, groupId uuid.UUID) ([]string, error) {
	query := "SELECT * from group_lesson_names($1)"
	result := []string{}
	err := repo.db.SelectContext(ctx, result, query, groupId)
	return result, err
}

func (repo *GroupsRepository) SelectSkippedHours(ctx context.Context, groupId uuid.UUID) ([]stored.SkippedHours, error) {
	query := "SELECT * from group_skipped_ordered()"
	result := []stored.SkippedHours{}
	err := repo.db.SelectContext(ctx, result, query)
	return result, err
}
