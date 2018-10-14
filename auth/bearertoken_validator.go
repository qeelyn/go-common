package auth

import (
	"context"
	"crypto/rsa"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"strings"
	"time"
)

type BearerTokenValidator struct {
	// Realm name to display to the user. Required.
	Realm string
	// Duration that a jwt token is valid. Optional, defaults to one hour.
	Timeout time.Duration
	// Callback function that should perform the authorization of the authenticated user. Called
	// only after an authentication success. Must return true on success, false on failure.
	// Optional, default to success.
	TokenValidator func(token *jwt.Token, c context.Context) error
	// Set the identity handler function. that mean the jwt is pass validete
	IdentityHandler func(c context.Context, claims jwt.MapClaims) (*Identity, error)
	// Secret EncryptionKey used for signing. Required.
	EncryptionKey []byte
	// Private EncryptionKey
	PrivKey *rsa.PrivateKey
	// Public EncryptionKey
	PubKey *rsa.PublicKey
}

var (
	// ErrMissingRealm indicates Realm name is required
	ErrMissingRealm = errors.New("realm is missing")

	// ErrForbidden when HTTP status 403 is given
	ErrForbidden = errors.New("you don't have permission to access this resource")

	// ErrInvalidPrivKey indicates that the given private EncryptionKey is invalid
	ErrInvalidClaims = errors.New("token payload content invalid")

	// ErrNoPrivKeyFile indicates that the given private EncryptionKey is unreadable
	ErrNoPrivKeyFile = errors.New("private EncryptionKey file unreadable")

	// ErrNoPubKeyFile indicates that the given public EncryptionKey is unreadable
	ErrNoPubKeyFile = errors.New("public EncryptionKey file unreadable")

	// ErrInvalidPrivKey indicates that the given private EncryptionKey is invalid
	ErrInvalidPrivKey = errors.New("private EncryptionKey invalid")

	// ErrInvalidPubKey indicates the the given public EncryptionKey is invalid
	ErrInvalidPubKey = errors.New("public EncryptionKey invalid")

	// ErrInvalidKey indicates the the given EncryptionKey is invalid
	ErrInvalidKey = errors.New("encrypty EncryptionKey invalid")
)

// pass through if private key is nil
func ParsePrivateKey(priKey []byte) (key *rsa.PrivateKey, err error) {
	if priKey == nil {
		return
	}
	key, err = jwt.ParseRSAPrivateKeyFromPEM(priKey)
	if err != nil {
		return
	}
	return
}

// pass through if public key is nil
func ParsePublicKey(pubKey []byte) (key *rsa.PublicKey, err error) {
	if pubKey == nil {
		return
	}
	key, err = jwt.ParseRSAPublicKeyFromPEM(pubKey)
	if err != nil {
		return
	}
	return
}

// Init initialize jwt configs.
func (b *BearerTokenValidator) Init() error {
	return nil
}

func (b *BearerTokenValidator) Validate(ctx context.Context, input string) (*Identity, error) {
	token, err := b.parseToken(input)

	if err != nil {
		return nil, err
	}
	// customer validate
	if b.TokenValidator != nil {
		if err = b.TokenValidator(token, ctx); err != nil {
			return nil, ErrForbidden
		}
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, ErrInvalidClaims
	}

	return b.IdentityHandler(ctx, claims)
}

func (b *BearerTokenValidator) parseToken(token string) (*jwt.Token, error) {
	return jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if strings.HasPrefix(token.Method.Alg(), "HS") {
			if len(b.EncryptionKey) == 0 {
				return nil, ErrInvalidKey
			}
			return b.EncryptionKey, nil
		} else {
			if b.PubKey == nil {
				return nil, ErrInvalidPubKey
			}
			return b.PubKey, nil
		}
	})
}
