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
	GetByEmail(ctx context.Context, universityId uuid.UUID, email, password string) (*stored.User, error)
	GetById(ctx context.Context, userId, universityId uuid.UUID) (*stored.User, error)
}
type AuthService struct {
	AuthRepository
}

func NewAuthService(repo AuthRepository) *AuthService {
	return &AuthService{AuthRepository: repo}
}

func (serv *AuthService) GetById(ctx context.Context, userId, universityId uuid.UUID) (*entities.User, error) {
	storedUsr, err := serv.AuthRepository.GetById(ctx, userId, universityId)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by ids from database in auth service: %w", err)
	}
	byteIds := []int8{}
	for _, role := range storedUsr.Roles {
		byteIds = append(byteIds, int8(role))
	}
	return entities.NewUser(storedUsr.Id, universityId, entities.RolesFromId(byteIds),
		storedUsr.Name, storedUsr.Surname, storedUsr.Patronymic,
		storedUsr.Email, storedUsr.Password), nil
}

func (serv *AuthService) GenerateTokenPair(ctx context.Context, email, password string, universityId uuid.UUID) (accessToken, refreshToken string, err error) {
	usr, err := serv.GetByEmail(ctx, universityId, email, password)
	if err != nil {
		return "", "", fmt.Errorf("failed to get user from database durin token pair generation: %w", err)
	}
	if usr == nil {
		return "", "", errors.ErrNoUsr
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
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	jwtAccess := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)

	refreshExpiration := time.Now().Add(7 * time.Hour * 24)
	refreshClaims := &secrets.Claims{
		UserID:       usrId,
		UniversityId: universityId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(refreshExpiration),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
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
	token, err := jwt.ParseWithClaims(refreshToken, claims, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Method.Alg())
		}
		return secrets.RefreshSalt, nil
	})
	if err != nil {
		return "", "", fmt.Errorf("failed to parse access token: %w", err)
	}
	return serv.getTokenPairs(ctx, token.Claims.(*secrets.Claims).UserID, token.Claims.(*secrets.Claims).UniversityId)
}

func (srv *AuthService) AddUser(ctx context.Context, input controllers.SignUpInput) error {
	err := srv.AuthRepository.SaveUser(ctx,
		entities.NewUser(uuid.UUID{}, input.UniversityId,
			entities.RolesFromId(input.UserRoles),
			input.Name, input.Surname, input.Patronymic,
			input.Email, input.Password))
	return err
}
