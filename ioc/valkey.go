package ioc

import (
	"db_labs/ioc/constants"
	"fmt"
	"log/slog"
	"net"
	"os"

	valkeyRepo "db_labs/repository/valkey"

	"github.com/valkey-io/valkey-go"
	"sigs.k8s.io/yaml"
)

type valkeyConf struct {
	Host string `json:"host"`
	Port int16  `json:"port"`
}

var UseValkeyConnection = provider(
	func() valkey.Client {
		fileData, err := os.ReadFile(constants.ValkeyConfPath)
		if err != nil {
			slog.Error(fmt.Errorf("Failed to read valkey conf at path %v: %w", constants.ValkeyConfPath, err).Error())
			os.Exit(1)
		}
		conf := valkeyConf{}
		err = yaml.Unmarshal(fileData, &conf)
		if err != nil {
			slog.Error(fmt.Errorf("Failed to unmarshal valkey yaml config at path %v: %w", constants.ValkeyConfPath, err).Error())
			os.Exit(1)
		}
		client, err := valkey.NewClient(valkey.ClientOption{
			InitAddress: []string{net.JoinHostPort(conf.Host, fmt.Sprint(conf.Port))},
		})
		if err != nil {
			slog.Error(fmt.Errorf("Failed to init valkey using config: %w", err).Error())
			os.Exit(-1)
		}
		return client
	},
)

var useTokenStore = provider(
	func() *valkeyRepo.AuthRepository {
		return valkeyRepo.NewAuthRepository(UseValkeyConnection())
	},
)
