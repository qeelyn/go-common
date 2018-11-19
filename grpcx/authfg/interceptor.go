package authfg

import (
	"context"
	"github.com/dgrijalva/jwt-go"
	"github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"github.com/grpc-ecosystem/go-grpc-middleware/util/metautils"
	"github.com/qeelyn/go-common/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const (
	authorizationKey  = "authorization"
	organizationIdKey = "orgid"
)

type authTag struct{}

var (
	authTagKey = &authTag{}
	ogrTagKey  = &authTag{}
)

// fromHttpHeader表示从http头部信息请求,一般为gateway时设置为true
func WithAuthClient(fromHttpHeader bool) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		var authHeader, orgid string
		if fromHttpHeader {
			if val := ctx.Value(authorizationKey); val != nil {
				authHeader = val.(string)
			}
			if val := ctx.Value(organizationIdKey); val != nil {
				orgid = val.(string)
			}
		} else {
			authHeader = metautils.ExtractIncoming(ctx).Get(authorizationKey)

			identity := ctx.Value(auth.ActiveUserContextKey)
			if identity != nil {
				if id, ok := identity.(auth.Identity); ok {
					orgid = id.OrgId
				}
			}
		}
		var kv []string
		if authHeader != "" {
			kv = append(kv, authorizationKey, authHeader)
		}
		if orgid != "" {
			kv = append(kv, organizationIdKey, orgid)
		}

		if len(kv) > 0 {
			ctx = metadata.AppendToOutgoingContext(ctx, kv...)
		}

		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

// auth base jwt,it parse BearToken to Identity entity
func ServerJwtAuthFunc(config map[string]interface{}) grpc_auth.AuthFunc {
	pubKey, _ := config["public-key"].([]byte)
	ekey, _ := config["encryption-key"].(string)
	publicKey, err := auth.ParsePublicKey(pubKey)
	if err != nil {
		panic(err)
	}
	validator := auth.BearerTokenValidator{
		PubKey:        publicKey,
		EncryptionKey: []byte(ekey),
		IdentityHandler: func(ctx context.Context, claims jwt.MapClaims) (*auth.Identity, error) {
			orgId := metautils.ExtractIncoming(ctx).Get(organizationIdKey)
			id, _ := claims["sub"].(string)
			identity := &auth.Identity{
				Id:    id,
				OrgId: orgId,
			}
			return identity, nil
		},
	}
	validator.Init()
	return func(ctx context.Context) (context.Context, error) {
		jwtHeaderPrefix := "bearer"
		token, err := grpc_auth.AuthFromMD(ctx, jwtHeaderPrefix)
		if err != nil {
			return ctx, err
		}
		if id, err := validator.Validate(ctx, token); err != nil {
			return nil, status.Error(codes.Unauthenticated, err.Error())
		} else {
			return context.WithValue(ctx, auth.ActiveUserContextKey, id), nil
		}
	}
}
