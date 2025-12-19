package controllers

import (
	"context"
	"db_labs/secrets"
	"db_labs/utils"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type AuthMiddleware struct {
	AuthService
}

func NewAuthMiddleware(serv AuthService) *AuthMiddleware {
	return &AuthMiddleware{AuthService: serv}
}

const (
	userIdCtxKey       = "usrid"
	universityIdCtxKey = "uniid"
)

func useAuthorized(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("access_token")
		if err != nil {
			utils.JSONError(w, "No access_token cookie", http.StatusForbidden)
			return
		}
		userId, universityId, err := parseToken(cookie.Value)
		if err != nil {
			utils.JSONError(w, "Invalid authorization token", http.StatusForbidden)
			return
		}
		r = r.WithContext(context.WithValue(r.Context(), universityIdCtxKey, universityId))
		r = r.WithContext(context.WithValue(r.Context(), userIdCtxKey, userId))
		next.ServeHTTP(w, r)
	})
}

func parseToken(accessToken string) (uuid.UUID, uuid.UUID, error) {
	claims := secrets.Claims{}
	token, err := jwt.ParseWithClaims(accessToken, &claims, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Method.Alg())
		}
		return secrets.AccessSalt, nil
	})
	if err != nil {
		return uuid.UUID{}, uuid.UUID{}, err
	}
	return token.Claims.(*secrets.Claims).UserID, token.Claims.(*secrets.Claims).UniversityId, err
}

func useTimeout(timeout time.Duration, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), timeout)
		defer cancel()
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
