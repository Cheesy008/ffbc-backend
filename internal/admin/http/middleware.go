package http

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"

	"github.com/cheesy008/ffbc-backend/internal/admin/constants"
	admincontext "github.com/cheesy008/ffbc-backend/internal/admin/http/context"
	"github.com/cheesy008/ffbc-backend/internal/admin/use_case"
)

func sessionMiddleware(authUseCase *use_case.AuthUseCase) func(huma.Context, func(huma.Context)) {
	return func(ctx huma.Context, next func(huma.Context)) {
		sessionToken, err := readSessionCookie(ctx.Header("Cookie"))
		if err != nil {
			writeUnauthorized(ctx)
			return
		}

		user, err := authUseCase.GetAdminBySessionToken(ctx.Context(), sessionToken)
		if err != nil {
			writeUnauthorized(ctx)
			return
		}

		next(huma.WithContext(ctx, admincontext.WithCurrentAdmin(ctx.Context(), user)))
	}
}

func readSessionCookie(cookieHeader string) (string, error) {
	req := http.Request{
		Header: http.Header{
			"Cookie": []string{cookieHeader},
		},
	}

	cookie, err := req.Cookie(constants.CookieName)
	if err != nil {
		return "", err
	}

	return cookie.Value, nil
}

func writeUnauthorized(ctx huma.Context) {
	ctx.SetStatus(http.StatusUnauthorized)
	ctx.SetHeader("Content-Type", "application/problem+json")
	_, _ = ctx.BodyWriter().Write([]byte(`{"title":"Unauthorized","status":401,"detail":"invalid admin session"}`))
}
