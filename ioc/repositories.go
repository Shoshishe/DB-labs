package ioc

import (
	"database/sql"
	"db_labs/ioc/constants"
	"fmt"
	_ "github.com/lib/pq"
	"log/slog"
	"os"
	"sigs.k8s.io/yaml"
)

type postgresConf struct {
	Host string `json:"host" binding:"required"`
	Port uint16 `json:"port" binding:"required"`
	Db   string `json:"db" binding:"required"`
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
		db, err := sql.Open("postgres", conf.Db)
		if err != nil {
			slog.Error(fmt.Errorf("Failed to connect to postgres database with given config: %w", err).Error())
			os.Exit(1)
		}
		return db
	},
)
