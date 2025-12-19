package postgres

import (
	"context"
	"database/sql"
	"db_labs/repository/postgres/stored"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

const (
	groupsTable = "groups"
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

func (repo *GroupsRepository) GetGroups(ctx context.Context, itemsPerPage uint8, currentPage uint) ([]stored.Group, error) {
	query := fmt.Sprintf("SELECT id, name, faculty_id FROM %s LIMIT $1 OFFSET $2", groupsTable)
	storedGroups := []stored.Group{}
	err := repo.db.SelectContext(ctx, &storedGroups, query, itemsPerPage, currentPage*uint(itemsPerPage))
	return storedGroups, err
}
