package auth

import (
	"context"
)

const (
	RoleAdmin  = "ADMIN"
	RoleAuthor = "AUTHOR"
)

type ctxKey int

const userKey ctxKey = 0
const rolesKey ctxKey = 1

func ContextWithUser(ctx context.Context, claims Claims) context.Context {
	ctx = context.WithValue(ctx, userKey, claims.UserID)
	ctx = context.WithValue(ctx, rolesKey, claims.Roles)
	return ctx
}

func UserFromContext(ctx context.Context) (string, bool) {
	id, ok := ctx.Value(userKey).(string)
	return id, ok
}

func RolesFromContext(ctx context.Context) ([]string, bool) {
	roles, ok := ctx.Value(rolesKey).([]string)
	return roles, ok
}

func HasRole(roles []string, wanted string) bool {
	for _, has := range roles {
		if has == wanted {
			return true
		}
	}
	return false
}
