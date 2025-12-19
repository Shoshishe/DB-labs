package postgres

import (
	"context"
	"database/sql"
	"db_labs/controllers"
	"db_labs/entities"
	"db_labs/repository/postgres/stored"
	"fmt"
	"strings"

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
	query := fmt.Sprintf("INSERT INTO %s (name, surname, patronymic, email, password) VALUES ($1,$2,$3,$4,$5) RETURNING id", usersTable)
	tx, err := repo.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback()
	row := tx.QueryRowContext(ctx, query, usr.Name(), usr.Surname(), usr.Patronymic(), usr.Email(), usr.Password())
	id := uuid.UUID{}
	err = row.Scan(&id)
	if err != nil {
		return err
	}

	if len(usr.Roles()) != 0 {
		for _, role := range usr.Roles() {
			query = fmt.Sprintf("INSERT INTO %s (user_id, role_id, university_id) VALUES ($1, $2, $3)", usersRolesTable)
			_, err = tx.ExecContext(ctx, query, id, role, usr.UniversityId())
			if err != nil {
				return err
			}
		}
	}
	err = tx.Commit()
	return err
}

func (repo *UserRepository) GetByEmail(ctx context.Context, universityId uuid.UUID, email, password string) (*stored.User, error) {
	query := fmt.Sprintf("SELECT u.id, u.name, u.surname, u.patronymic, u.email, ARRAY(SELECT role_id FROM %s as ur WHERE ur.university_id=$3 AND ur.user_id=u.id) as roles FROM %s AS u where u.email=$1 AND u.password=$2", usersRolesTable, usersTable)
	usr := &stored.User{}
	err := repo.db.GetContext(ctx, usr, query, email, password, universityId)
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

func (repo *UserRepository) GetRoles(ctx context.Context) ([]int8, error) {
	query := fmt.Sprintf("SELECT id FROM %s", rolesTable)
	result := []int8{}
	err := repo.db.SelectContext(ctx, &result, query)
	return result, err
}

func (repo *UserRepository) UpdateUser(ctx context.Context, request controllers.UpdateUserRequest) error {
	argPosition := 1
	args := []any{}
	positionedArgs := []string{}
	if request.Name != nil {
		positionedArgs = append(positionedArgs, fmt.Sprintf("name=$%d", argPosition))
		args = append(args, request.Name)
		argPosition++
	}
	if request.Patronymic != nil {
		positionedArgs = append(positionedArgs, fmt.Sprintf("patronymic=$%d", argPosition))
		args = append(args, request.Patronymic)
		argPosition++
	}
	if request.Surname != nil {
		positionedArgs = append(positionedArgs, fmt.Sprintf("surname=$%d", argPosition))
		args = append(args, request.Surname)
		argPosition++
	}
	if argPosition == 1 {
		return nil
	}
	setQuery := strings.Join(positionedArgs, ", ")
	args = append(args, request.Id)
	query := fmt.Sprintf("UPDATE %s SET %s WHERE id=$%d", usersTable, setQuery, argPosition)
	if request.Roles == nil {
		_, err := repo.db.ExecContext(ctx, query, args...)
		return err
	} else {
		tx, err := repo.db.BeginTx(ctx, &sql.TxOptions{})
		if err != nil {
			return err
		}
		defer tx.Rollback()
		_, err = tx.ExecContext(ctx, query, args...)
		if err != nil {
			return err
		}
		query = fmt.Sprintf("DELETE FROM %s WHERE user_id=$1", usersRolesTable)
		_, err = tx.ExecContext(ctx, query, request.Id)
		if err != nil {
			return err
		}
		query = fmt.Sprintf("INSERT INTO %s user_id, role_id, university_id VALUES ($1, $2, $3)", usersRolesTable)
		for _, role := range *request.Roles {
			_, err = tx.ExecContext(ctx, query, request.Id, role, request.UniversityId)
			if err != nil {
				return err
			}
		}
		return tx.Commit()
	}
}
