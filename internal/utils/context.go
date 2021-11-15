package utils

import (
	"context"

	"github.com/google/uuid"
)

type CtxKey int

const (
	UserID CtxKey = iota
	ContentType
)

func ContextWithUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, UserID, userID)
}

func UserIDFromContext(ctx context.Context) (uuid.UUID, bool) {
	token, ok := ctx.Value(UserID).(string)
	if !ok {
		return uuid.UUID{}, false
	}

	userID, err := uuid.Parse(token)
	if err != nil {
		return uuid.UUID{}, false

	}
	return userID, true
}

func ContextWithContentType(ctx context.Context, contentType string) context.Context {
	return context.WithValue(ctx, ContentType, contentType)
}

func ContentTypeFromContext(ctx context.Context) (string, bool) {
	contentType, ok := ctx.Value(ContentType).(string)
	return contentType, ok
}
