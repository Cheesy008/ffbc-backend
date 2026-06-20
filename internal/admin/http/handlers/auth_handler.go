package handlers

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/danielgtaylor/huma/v2"

	"github.com/cheesy008/ffbc-backend/internal/admin/constants"
	"github.com/cheesy008/ffbc-backend/internal/admin/domain"
	"github.com/cheesy008/ffbc-backend/internal/admin/use_case"
)

type AuthHandler struct {
	authUseCase *use_case.AuthUseCase
}

func NewAuthHandler(authUseCase *use_case.AuthUseCase) *AuthHandler {
	return &AuthHandler{authUseCase: authUseCase}
}

type LoginInput struct {
	Body LoginRequest
}

type LoginRequest struct {
	Email    string `json:"email" format:"email" required:"true" example:"admin@example.com"`
	Password string `json:"password" required:"true" example:"password"`
}

type LoginOutput struct {
	SetCookie string `header:"Set-Cookie"`
	Body      LoginResponse
}

type LoginResponse struct {
	Message string `json:"message" example:"string"`
}

type LogoutInput struct {
	SessionToken http.Cookie `cookie:"admin_session"`
}

type LogoutOutput struct {
	SetCookie string   `header:"Set-Cookie"`
	_         struct{} `status:"200"`
}

func (h *AuthHandler) Login(ctx context.Context, input *LoginInput) (*LoginOutput, error) {
	output, err := h.authUseCase.Login(ctx, use_case.LoginInput{
		Email:         input.Body.Email,
		PlainPassword: input.Body.Password,
	})
	if err != nil {
		return nil, mapError(err)
	}

	return &LoginOutput{
		SetCookie: sessionCookie(output.SessionToken, constants.CookieExpirationTime),
		Body: LoginResponse{
			Message: "ok",
		},
	}, nil
}

func sessionCookie(sessionToken string, ttl time.Duration) string {
	cookie := http.Cookie{
		Name:     constants.CookieName,
		Value:    sessionToken,
		Path:     "/api/admin",
		MaxAge:   int(ttl.Seconds()),
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}

	return cookie.String()
}

func expiredSessionCookie() string {
	cookie := http.Cookie{
		Name:     constants.CookieName,
		Value:    "",
		Path:     "/api/admin",
		MaxAge:   -1,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}

	return cookie.String()
}

func (h *AuthHandler) Logout(ctx context.Context, input *LogoutInput) (*LogoutOutput, error) {
	err := h.authUseCase.Logout(ctx, use_case.LogoutInput{SessionToken: input.SessionToken.Value})
	if err != nil {
		return nil, mapError(err)
	}

	return &LogoutOutput{
		SetCookie: expiredSessionCookie(),
	}, nil
}

func mapError(err error) error {
	switch {
	case errors.Is(err, domain.ErrAdminNotFound), errors.Is(err, domain.ErrInvalidCredentials):
		return huma.Error401Unauthorized("invalid credentials")
	case errors.Is(err, domain.ErrAdminInactive):
		return huma.Error403Forbidden("admin inactive")
	default:
		return huma.Error500InternalServerError("internal server error")
	}
}
