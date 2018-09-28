package auth

import (
	"context"
	"errors"
	"strconv"
)

type contextKey struct{}

var ActiveUserContextKey = contextKey{}

type Identity struct {
	// user id
	Id string
	// org id
	OrgId string
}

// 获取int格式的ID,如果id为int的话
func (t *Identity) IdInt() int32 {
	i, _ := strconv.Atoi(t.Id)
	return int32(i)
}

// 获取int格式的ID,如果id为int的话
func (t *Identity) OrgIdInt() int32 {
	i, _ := strconv.Atoi(t.OrgId)
	return int32(i)
}

// get User Id from context, grpc interceptor convert metadata into context
func UserFromContext(ctx context.Context) (*Identity, error) {
	if user, ok := ctx.Value(ActiveUserContextKey).(*Identity); ok {
		if user.Id != "" {
			return user, nil
		}
	}
	return nil, errors.New("can't find user id")
}
