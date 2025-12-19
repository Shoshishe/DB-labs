package ioc

import (
	"database/sql"
	"db_labs/ioc/constants"
	repository "db_labs/repository/postgres"
	"fmt"
	"log/slog"
	"os"

	_ "github.com/lib/pq"
	"sigs.k8s.io/yaml"
)

type postgresConf struct {
	Host string `json:"host" binding:"required"`
	Port uint16 `json:"port" binding:"required"`
	Db   string `json:"db" binding:"required"`
	User string `json:"user" binding:"required"`
}

var UsePgConnection = provider(
	func() *sql.DB {
		fileData, err := os.ReadFile(constants.PostgresConfPath)
		if err != nil {
			slog.Error(fmt.Errorf("Failed to read postgres conf at path %v: %w", constants.PostgresConfPath, err).Error())
			os.Exit(1)
		}
		conf := postgresConf{}
		err = yaml.Unmarshal(fileData, &conf)
		if err != nil {
			slog.Error(fmt.Errorf("Failed to unmarshal yaml config at path %v: %w", constants.PostgresConfPath, err).Error())
			os.Exit(1)
		}
		db, err := sql.Open("postgres", fmt.Sprintf("user=%s host=%s port=%d password=%s dbname=%s sslmode=disable", conf.User, conf.Host, conf.Port, os.Getenv("DB_PASSWORD"), conf.Db))
		if err != nil {
			slog.Error(fmt.Errorf("Failed to connect to postgres database with given config: %w", err).Error())
			os.Exit(1)
		}
		err = db.Ping()
		db.SetMaxOpenConns(10)
		if err != nil {
			slog.Error(fmt.Errorf("Failed to connect to postgres database with given config: %w", err).Error())
			os.Exit(1)
		}
		return db
	},
)

var UseUniversitiesRepo = provider(
	func() *repository.UniversitiesRepository {
		return repository.NewUniversitiesRepository(UsePgConnection())
	},
)

var useFacultiesRepo = provider(
	func() repository.FacultiesRepository {
		return *repository.NewFacultiesRepository(UsePgConnection())
	},
)

var useLessonsRepo = provider(
	func() repository.LessonsRepository {
		return *repository.NewLessonsRepository(UsePgConnection())
	},
)

var useUsersRepo = provider(
	func() *repository.UserRepository {
		return repository.NewUserRepository(UsePgConnection())
	},
)

var useGroupsRepo = provider(
	func() *repository.GroupsRepository {
		return repository.NewGroupsRepository(UsePgConnection())
	},
)
