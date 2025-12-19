package valkey

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/valkey-io/valkey-go"
)

type AuthRepository struct {
	db valkey.Client
}

func NewAuthRepository(db valkey.Client) *AuthRepository {
	return &AuthRepository{db: db}
}

func (repo *AuthRepository) SaveToken(ctx context.Context, userId uuid.UUID, refreshToken string) error {
	res := repo.db.Do(ctx, repo.db.B().Set().Key(userId.String()).Value(refreshToken).Ex(7*24*time.Hour).Build())
	return res.Error()
}
