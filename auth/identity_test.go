package auth_test

import (
	"context"
	"github.com/qeelyn/go-common/auth"
	"testing"
)

func TestUserFromContext(t *testing.T) {
	ctx := context.Background()
	ctx = context.WithValue(ctx, auth.ActiveUserContextKey, &auth.Identity{Id: "12", OrgId: "13"})
	uc, err := auth.UserFromContext(ctx)
	if err != nil {
		t.Error(err)
	}
	if uc.IdInt() != 12 || uc.OrgIdInt() != 13 {
		t.Error()
	}
}
