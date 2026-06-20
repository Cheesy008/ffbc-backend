package context

import (
	stdcontext "context"
	"fmt"

	"github.com/cheesy008/ffbc-backend/internal/admin/domain"
)

type currentAdminContextKey struct{}

func WithCurrentAdmin(ctx stdcontext.Context, admin domain.AdminUser) stdcontext.Context {
	return stdcontext.WithValue(ctx, currentAdminContextKey{}, admin)
}

func CurrentAdmin(ctx stdcontext.Context) (domain.AdminUser, bool) {
	admin, ok := ctx.Value(currentAdminContextKey{}).(domain.AdminUser)
	return admin, ok
}

func MustCurrentAdmin(ctx stdcontext.Context) domain.AdminUser {
	admin, ok := CurrentAdmin(ctx)
	if !ok {
		panic(fmt.Errorf("current admin is missing from protected route context"))
	}

	return admin
}
