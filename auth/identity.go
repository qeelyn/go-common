package auth

import (
	"context"
	"errors"
)

type contextKey struct{}

var ActiveUserContextKey = contextKey{}

type Identity struct {
	// user id
	Id int32
	// org id
	OrgId int32
}

type UserContext struct {
	UserId int32
	OrgId  int32
}

// get User Id from context, grpc interceptor convert metadata into context
func UserFromContext(ctx context.Context) (*UserContext, error) {
	ua := &UserContext{}
	if userCtx, ok := ctx.Value(ActiveUserContextKey).(*Identity); ok {
		if userCtx.Id > 0 {
			ua.UserId = userCtx.Id
		} else {
			return nil, errors.New("can't find user id")
		}
		// no org is ok
		ua.OrgId = userCtx.OrgId
	}
	return ua, nil
}
