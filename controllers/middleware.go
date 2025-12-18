package controllers

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"
)

func useCors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type,access-control-allow-origin, access-control-allow-headers")
		next.ServeHTTP(w, r)
	})
}

type AuthMiddleware struct {
	AuthService
}

const userIdCtxKey = "usrid"

func (mdw *AuthMiddleware) useAuthorized(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		if header == "" {
			http.Error(w, "Empty authorization header", http.StatusForbidden)
			return
		}
		headerParts := strings.Split(header, " ")
		if len(headerParts) != 2 {
			http.Error(w, fmt.Sprintf("Invalid count of header parts. Expected two, got %d", len(headerParts)), http.StatusForbidden)
			return
		}
		userId, err := mdw.ParseToken(headerParts[1])
		if err != nil {
			http.Error(w, "Invalid authorization token", http.StatusForbidden)
			return
		}

		r.WithContext(context.WithValue(r.Context(), userIdCtxKey, userId))
		next.ServeHTTP(w, r)
	})
}

func useTimeout(timeout time.Duration, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), timeout)
		defer cancel()
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
