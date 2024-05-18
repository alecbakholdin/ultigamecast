package ctx_var

import (
	"context"
	"fmt"
	"ultigamecast/models"
)

type ContextVar string

const (
	HttpMethod ContextVar = "method"
	Path       ContextVar = "path"
	ReqId      ContextVar = "RequestId"
	User       ContextVar = "User"
)

var LogMessageVars = []ContextVar{HttpMethod, Path}
var LogAttrVars = []ContextVar{ReqId, User}

func IsAuthenticated(ctx context.Context) bool {
	return GetUser(ctx) != nil
}

func GetUser(ctx context.Context) *models.User {
	if user, ok := ctx.Value(User).(*models.User); ok {
		return user
	}
	return nil
}

func GetValue(ctx context.Context, key ContextVar) string {
	if key == User {
		if user := GetUser(ctx); user != nil {
			return fmt.Sprintf("[%d] %s", user.ID, user.Email)
		}
	} else if val, ok := ctx.Value(key).(string); ok {
		return val
	}
	return ""
}
