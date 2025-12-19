package controllers

import (
	"context"
	"crypto/sha256"
	"db_labs/controllers/constants"
	"db_labs/controllers/responses"
	"db_labs/entities"
	"db_labs/secrets"
	"db_labs/utils"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/google/uuid"
)

const (
	refreshTimeout = 7 * 24 * time.Hour
	accessTimeout  = 15 * time.Minute
)

type AuthService interface {
	GenerateTokenPair(ctx context.Context, email string, password string, universityId uuid.UUID) (accessToken, refreshToken string, err error)
	// ParseToken(accessToken string) (userId, universityId uuid.UUID, err error)
	AddUser(context.Context, SignUpInput) error
	RegenerateTokens(ctx context.Context, refreshToken string) (refresh string, access string, err error)
	GetById(ctx context.Context, usrId, uniId uuid.UUID) (*entities.User, error)
}

type AuthController struct {
	serv AuthService
}

func NewAuthController(mux *http.ServeMux, serv AuthService) *AuthController {
	controller := &AuthController{serv: serv}
	return controller
}

func (controller *AuthController) RegisterRoutes(mux *http.ServeMux) {
	mux.Handle("POST /auth/sign-in", useTimeout(constants.DefaultTimeout, (http.HandlerFunc(controller.signIn))))
	mux.Handle("POST /auth/sign-up", useTimeout(constants.DefaultTimeout, (http.HandlerFunc(controller.signUp))))
	mux.Handle("POST /auth/refresh", useTimeout(constants.DefaultTimeout, (http.HandlerFunc(controller.refresh))))
	mux.Handle("POST /auth/logout", useTimeout(constants.DefaultTimeout, (http.HandlerFunc(controller.logout))))
	mux.Handle("GET /auth/me", useTimeout(constants.DefaultTimeout, useAuthorized(http.HandlerFunc(controller.me))))
}

type SignUpInput struct {
	Email        string    `json:"email"`
	Password     string    `json:"password"`
	UniversityId uuid.UUID `json:"university_id"`
	UserRoles    []int8    `json:"roles"`
	Patronymic   string    `json:"patronymic"`
	Name         string    `json:"name"`
	Surname      string    `json:"surname"`
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
	if err != nil {
		utils.JSONError(wr, fmt.Sprintf("Failed to add user: %v", err), http.StatusForbidden)
	}
}

type SignInInput struct {
	Email        string    `json:"email" binding:"required"`
	Password     string    `json:"password" binding:"required"`
	UniversityId uuid.UUID `json:"university_id" binding:"required"`
}

func (controller *AuthController) signIn(wr http.ResponseWriter, req *http.Request) {
	var input SignInInput
	if err := json.NewDecoder(req.Body).Decode(&input); err != nil {
		slog.Error("Failed to decode sign in input", "err", err.Error())
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
		slog.Error("Failed to sign in user", "err", err.Error())
		utils.JSONError(wr, err.Error(), http.StatusBadRequest)
		return
	}
	controller.setTokensCookie(&wr, access, refresh)

	wr.WriteHeader(http.StatusOK)
	err = json.NewEncoder(wr).Encode("ok")
	if err != nil {
		slog.Error("Failed to send ok response", "err", err)
	}
}

func (controller *AuthController) setTokensCookie(wr *http.ResponseWriter, accessToken, refreshToken string) {
	http.SetCookie(*wr, &http.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		Expires:  time.Now().Add(accessTimeout),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
	})

	http.SetCookie(*wr, &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Expires:  time.Now().Add(refreshTimeout),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
	})
}

func (controller *AuthController) refresh(wr http.ResponseWriter, req *http.Request) {
	refershCookie, err := req.Cookie("refresh_token")
	if err != nil {
		slog.Error("Failed to get refresh token cookie: %s", slog.String("err", err.Error()))
		utils.JSONError(wr, "Refresh token not found", http.StatusUnauthorized)
		return
	}
	access, refresh, err := controller.serv.RegenerateTokens(req.Context(), refershCookie.Value)
	if err != nil {
		slog.Error("Failed to regenerate access and refresh token pair: %s", slog.String("err", err.Error()))
		utils.JSONError(wr, "Couldn't regenerate access and refresh tokens", http.StatusUnauthorized)
		return
	}
	controller.setTokensCookie(&wr, access, refresh)
	wr.WriteHeader(http.StatusOK)
	fmt.Fprintf(wr, "Cookie 'session_token' has been set!")
}

func (controller *AuthController) logout(wr http.ResponseWriter, _ *http.Request) {
	deleteCookie("access_token", wr, true, false, http.SameSiteLaxMode)
	deleteCookie("refresh_token", wr, true, false, http.SameSiteLaxMode)
	err := json.NewEncoder(wr).Encode("ok")
	if err != nil {
		slog.Error("Failed to send ok response", "err", err)
	}
}

func (controller *AuthController) me(wr http.ResponseWriter, r *http.Request) {
	userId, ok := r.Context().Value(userIdCtxKey).(uuid.UUID)
	if !ok {
		utils.JSONError(wr, "Failed to get credentials", http.StatusBadRequest)
		return
	}

	universityId, ok := r.Context().Value(userIdCtxKey).(uuid.UUID)
	if !ok {
		utils.JSONError(wr, "Failed to get credentials", http.StatusBadRequest)
		return
	}

	usr, err := controller.serv.GetById(r.Context(), userId, universityId)
	if err != nil {
		utils.JSONError(wr, fmt.Sprintf("Failed to get user: %v", err), http.StatusInternalServerError)
		return
	}

	resp := responses.MeFromEntity(usr)
	err = json.NewEncoder(wr).Encode(&resp)
	if err != nil {
		slog.Error("Failed to marshal resposne: %w", "err", err)
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
