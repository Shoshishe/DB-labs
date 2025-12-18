package services

import (
	"context"
	"db_labs/controllers"
	"db_labs/entities"
	"db_labs/repository/postgres/stored"
	"db_labs/secrets"
	"db_labs/services/errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type AuthRepository interface {
	SaveToken(ctx context.Context, userId uuid.UUID, refershToken string) error
	SaveUser(ctx context.Context, user *entities.User) error
	GetByEmail(ctx context.Context, enmail, password string) (*stored.User, error)
}
type AuthService struct {
	AuthRepository
}

func NewAuthService(repo AuthRepository) *AuthService {
	return &AuthService{AuthRepository: repo}
}

func (serv AuthService) GenerateTokenPair(ctx context.Context, email, password string, universityId uuid.UUID) (accessToken, refreshToken string, err error) {
	usr, err := serv.GetByEmail(ctx, email, password)
	if usr == nil {
		return "", "", errors.ErrNoUsr
	}
	if err != nil {
		return "", "", fmt.Errorf("failed to get user from database durin token pair generation: %w", err)
	}
	return serv.getTokenPairs(ctx, usr.Id, universityId)
}

func (serv *AuthService) getTokenPairs(ctx context.Context, usrId, universityId uuid.UUID) (accessToken, refreshToken string, err error) {
	accessExpiration := time.Now().Add(15 * time.Minute)
	accessClaims := &secrets.Claims{
		UserID:       usrId,
		UniversityId: universityId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(accessExpiration),
		},
	}
	jwtAccess := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)

	refreshExpiration := time.Now().Add(7 * time.Hour * 24)
	refreshClaims := &secrets.Claims{
		UserID:       usrId,
		UniversityId: universityId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(refreshExpiration),
		},
	}
	jwtRefresh := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)

	accessToken, err = jwtAccess.SignedString(secrets.AccessSalt)
	if err != nil {
		return "", "", err
	}
	refreshToken, err = jwtRefresh.SignedString(secrets.RefreshSalt)
	if err != nil {
		return "", "", err
	}

	err = serv.AuthRepository.SaveToken(ctx, usrId, refreshToken)
	if err != nil {
		return "", "", err
	}
	return accessToken, refreshToken, err
}

func (serv *AuthService) RegenerateTokens(ctx context.Context, refreshToken string) (refresh string, access string, err error) {
	claims := &secrets.Claims{}
	token, err := jwt.ParseWithClaims(refreshToken, claims.RegisteredClaims, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Method.Alg())
		}
		return secrets.AccessSalt, nil
	})
	if err != nil {
		return "", "", fmt.Errorf("failed to parse access token: %w", err)
	}
	return serv.getTokenPairs(ctx, token.Claims.(secrets.Claims).UserID, token.Claims.(secrets.Claims).UniversityId)
}

func (srv *AuthService) ParseToken(accessToken string) (uuid.UUID, error) {
	claims := secrets.Claims{}
	token, err := jwt.ParseWithClaims(accessToken, claims.RegisteredClaims, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Method.Alg())
		}
		return secrets.AccessSalt, nil
	})
	if err != nil {
		return uuid.UUID{}, err
	}
	return token.Claims.(secrets.Claims).UserID, err
}

func (srv *AuthService) AddUser(ctx context.Context, input controllers.SignUpInput) error {
	err := srv.AuthRepository.SaveUser(ctx,
		entities.NewUser(uuid.UUID{}, input.UniversityId,
			entities.RolesFromId(input.UserRoles),
			input.Name, input.Surname, input.Patronymic,
			input.Email, input.Password))
	return err
}
