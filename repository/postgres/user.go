package postgres

import (
	"context"
	"database/sql"
	"db_labs/entities"
	"db_labs/repository/postgres/stored"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

const (
	usersTable      = "users"
	usersRolesTable = "user_roles"
	rolesTable      = "roles"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: sqlx.NewDb(db, "pq")}
}

func (repo *UserRepository) SaveUser(ctx context.Context, usr *entities.User) error {
	query := fmt.Sprintf("INSERT INTO %s (name, surname, patronymic, email, password) VALUES ($1,$2,$3,$4,$5)", usersTable)
	tx, err := repo.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback()
	_, err = tx.ExecContext(ctx, query, usr.Name(), usr.Surname(), usr.Patronymic(), usr.Email(), usr.Password())
	if err != nil {
		return err
	}

	if len(usr.Roles()) != 0 {
		values := []any{}
		query = fmt.Sprintf("INSERT INTO %s (user_id, role_id, university_id) VALUES ", usersRolesTable)
		for i, role := range usr.Roles() {
			query += "(?, ?, ?)"
			if i < len(usr.Roles())-1 {
				query += ","
			}
			values = append(values, usr.Id(), role, usr.UniversityId())
		}
		_, err = tx.ExecContext(ctx, query, values...)
		err = tx.Commit()
	}
	return err
}

func (repo *UserRepository) GetByEmail(ctx context.Context, email, password string) (*stored.User, error) {
	query := fmt.Sprintf("SELECT id, name, surname. patronymic, email FROM %s where email=$1 AND password=$2", usersTable)
	usr := &stored.User{}
	err := repo.db.SelectContext(ctx, &usr, query, email, password)
	if err != nil {
		return nil, err
	}
	return usr, nil
}

func (repo *UserRepository) GetById(ctx context.Context, id, universityId uuid.UUID) (*stored.User, error) {
	tx, err := repo.db.BeginTxx(ctx, &sql.TxOptions{})
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	query := fmt.Sprintf("select id,name,surname,patronymic, email from %s where id=$1", usersTable)
	user := &stored.User{}
	err = tx.GetContext(ctx, user, query, id)
	if err != nil {
		return nil, err
	}
	roles := []int8{}
	query = fmt.Sprintf("SELECT r.id from %s AS ur WHERE ur.user_id=$1 AND ur.university_id=$2 INNER JOIN %s AS r ON r.id=ur.role_id", usersRolesTable, rolesTable)
	repo.db.SelectContext(ctx, &roles, query, id, universityId)
	return user, nil
}

// func (repo *UserRepository) UpdateUser(ctx context.Context, usr *entities.User) error {

// }
