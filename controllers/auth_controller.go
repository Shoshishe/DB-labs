package controllers

import (
	"context"
	"crypto/sha256"
	"db_labs/secrets"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type AuthService interface {
	GenerateTokenPair(ctx context.Context, email string, password string, universityId uuid.UUID) (accessToken, refreshToken string, err error)
	ParseToken(accessToken string) (userId uuid.UUID, err error)
	AddUser(context.Context, SignUpInput) error
	RegenerateTokens(ctx context.Context, refreshToken string) (refresh string, access string, err error)
}

type AuthController struct {
	serv AuthService
}

func NewAuthController(mux *http.ServeMux, serv AuthService) *AuthController {
	controller := &AuthController{serv: serv}
	controller.RegisterRoutes(mux)
	return controller
}

func (controller *AuthController) RegisterRoutes(mux *http.ServeMux) {
	mux.Handle("POST /auth/sign-in", useCors(http.HandlerFunc(controller.signIn)))
	mux.Handle("POST /auth/sign-up", useCors(http.HandlerFunc(controller.signUp)))
	mux.Handle("POST /auth/refresh", useCors(http.HandlerFunc(controller.refresh)))
	mux.Handle("POST /auth/logout", useCors(http.HandlerFunc(controller.logout)))
}

type SignUpInput struct {
	Email        string    `json:"email" binding:"required"`
	Password     string    `json:"password" binding:"required"`
	UniversityId uuid.UUID `json:"uni_id" binding:"required"`
	UserRoles    []int8    `json:"roles" binding:"required"`
	Patronymic   string    `json:"patronymic" binding:"required"`
	Name         string    `json:"name" binding:"required"`
	Surname      string    `json:"surname" binding:"required"`
}

func (controller *AuthController) signUp(wr http.ResponseWriter, r *http.Request) {
	input := SignUpInput{}
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		wr.WriteHeader(http.StatusBadRequest)
		res, _ := json.Marshal(fmt.Errorf("Failed to decode sign up input: %w", err).Error())
		wr.Write(res)
		return
	}
	hash := sha256.New()
	hash.Write([]byte(input.Password))
	input.Password = fmt.Sprintf("%x", hash.Sum([]byte(secrets.PasswordSalt)))
	err = controller.serv.AddUser(r.Context(), input)
}

type SignInInput struct {
	Email        string    `json:"email" binding:"required"`
	Password     string    `json:"password" binding:"required"`
	UniversityId uuid.UUID `json:"uni_id" binding:"required"`
}

func (controller *AuthController) signIn(wr http.ResponseWriter, req *http.Request) {
	var input SignInInput
	if err := json.NewDecoder(req.Body).Decode(&input); err != nil {
		slog.Error("Failed to decode sign in input" + err.Error())
		wr.WriteHeader(http.StatusBadRequest)
		res, _ := json.Marshal(err.Error)
		wr.Write(res)
		return
	}
	hash := sha256.New()
	hash.Write([]byte(input.Password))
	input.Password = fmt.Sprintf("%x", hash.Sum([]byte(secrets.PasswordSalt)))

	access, refresh, err := controller.serv.GenerateTokenPair(req.Context(), input.Email, input.Password, input.UniversityId)
	if err != nil {
		slog.Error(err.Error())
		wr.WriteHeader(http.StatusBadRequest)
		res, _ := json.Marshal(err.Error())
		wr.Write(res)
		return
	}
	controller.setTokensCookie(wr, access, refresh)
	err = json.NewEncoder(wr).Encode("ok")
	if err != nil {
		slog.Error("Failed to send ok response", "err", err)
	}
}

func (controller *AuthController) setTokensCookie(wr http.ResponseWriter, accessToken, refreshToken string) {
	http.SetCookie(wr, &http.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		Expires:  time.Now().Add(15 * time.Minute),
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
	})

	http.SetCookie(wr, &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Expires:  time.Now().Add(7 * 24 * time.Hour),
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
		Path:     "/auth/refresh",
	})
}

func (controller *AuthController) refresh(wr http.ResponseWriter, req *http.Request) {
	refershCookie, err := req.Cookie("refresh_token")
	if err != nil {
		slog.Error("Failed to get refresh token cookie: %s", slog.String("err", err.Error()))
		http.Error(wr, "Refresh token not found", http.StatusUnauthorized)
		return
	}

	access, refresh, err := controller.serv.RegenerateTokens(req.Context(), refershCookie.Value)
	if err != nil {
		slog.Error("Failed to regenerate access and refresh token pair: %s", slog.String("err", err.Error()))
		http.Error(wr, "Couldn't regenerate access and refresh tokens", http.StatusUnauthorized)
		return
	}
	controller.setTokensCookie(wr, access, refresh)
}

func (controller *AuthController) logout(wr http.ResponseWriter, _ *http.Request) {
	deleteCookie("access_token", wr, true, false, http.SameSiteLaxMode)
	deleteCookie("refresh_token", wr, true, false, http.SameSiteLaxMode)
	err := json.NewEncoder(wr).Encode("ok")
	if err != nil {
		slog.Error("Failed to send ok response", "err", err)
	}
}

func deleteCookie(name string, w http.ResponseWriter, HttpOnly, Secure bool, SameSite http.SameSite) {
	cookie := &http.Cookie{
		Name:     name,
		Value:    "",
		MaxAge:   -1,
		Expires:  time.Now().Add(-time.Hour),
		HttpOnly: HttpOnly,
		Secure:   Secure,
		SameSite: SameSite,
	}
	http.SetCookie(w, cookie)
}
