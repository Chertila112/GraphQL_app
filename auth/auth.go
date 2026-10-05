package auth

import "context"

type ctxKey string

const userIDKey ctxKey = "userID"

// WithUserID кладёт ID пользователя в контекст.
func WithUserID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, userIDKey, id)
}

// UserIDFromContext достаёт ID пользователя из контекста.
func UserIDFromContext(ctx context.Context) (string, bool) {
	id, ok := ctx.Value(userIDKey).(string)
	return id, ok
}