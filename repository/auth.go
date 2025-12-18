package repository

import (
	"db_labs/repository/postgres"
	"db_labs/repository/valkey"
)

type AuthRepistory struct {
	valkey.AuthRepository
	postgres.UserRepository
}

func NewAuthRepository(tokenStore valkey.AuthRepository, userRepo postgres.UserRepository) *AuthRepistory {
	return &AuthRepistory{AuthRepository: tokenStore, UserRepository: userRepo}
}
