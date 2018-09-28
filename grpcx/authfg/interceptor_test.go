package authfg_test

import (
	"context"
	"errors"
	"github.com/grpc-ecosystem/go-grpc-middleware/util/metautils"
	"github.com/qeelyn/go-common/auth"
	"github.com/qeelyn/go-common/grpcx/authfg"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"testing"
)

func TestWithAuthClient(t *testing.T) {
	ctx := context.TODO()
	header := "Bearer test"
	var hasAuth, hasOrg = false, false
	invoker := func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, opts ...grpc.CallOption) error {
		amd := metautils.ExtractOutgoing(ctx)
		if hasAuth {
			if h := amd.Get("authorization"); h != header {
				return errors.New("auth no eq")
			}

		}
		if hasOrg {
			if identity := ctx.Value(auth.ActiveUserContextKey); identity == nil {
				return errors.New("identity no found")
			} else {
				if i := identity.(auth.Identity); i.IdInt() != 22 && i.OrgIdInt() != 0 {
					return errors.New("identity value error")
				}
			}
		}
		return nil
	}
	// http mock
	uci := authfg.WithAuthClient(true)
	ctx = context.WithValue(ctx, "authorization", header)
	hasAuth = true
	err := uci(ctx, "", nil, nil, nil, invoker)
	if err != nil {
		t.Fatal(err)
	}

	// grpc mock
	ctx = context.WithValue(ctx, auth.ActiveUserContextKey, auth.Identity{
		Id:    "22",
		OrgId: "0",
	})
	ctx = metadata.NewIncomingContext(ctx, metadata.Pairs("authorization", header))
	hasAuth = true
	hasOrg = true
	err = uci(ctx, "", nil, nil, nil, invoker)
	if err != nil {
		t.Fatal(err)
	}
}

func TestServerJwtAuthFunc(t *testing.T) {
	cnf := map[string]interface{}{
		"public-key":     "",
		"encryption-key": "123456",
		"algorithm":      "HS256",
	}
	afunc := authfg.ServerJwtAuthFunc(cnf)

	badHeader := "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyLCJleHAiOjE1MTYyMzkxMjJ9.JcRoPW5fA44i7vuGyXGXKHuAfZYly_uFGs5FznyPJBc"
	okHeader := "bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.keH6T3x1z7mmhKL1T3r9sQdAxxdzB6siemGMr_6ZOwU"

	ctx := metadata.NewIncomingContext(context.TODO(), metadata.Pairs("authorization", okHeader, "orgid", "100"))
	newCtx, err := afunc(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if id := newCtx.Value(auth.ActiveUserContextKey); id == nil {
		t.Fatal("id no found")
	} else if idd, ok := id.(*auth.Identity); !ok || idd.OrgIdInt() != 100 {
		t.Fatal("org no found")
	}

	ctx = metadata.NewIncomingContext(context.TODO(), metadata.Pairs("authorization", badHeader, "orgid", "100"))
	_, err = afunc(ctx)
	if err == nil {
		t.Fatal("must be error")
	}
}
