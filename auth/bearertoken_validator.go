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
	// Secret key used for signing. Required.
	Key []byte

	// Duration that a jwt token is valid. Optional, defaults to one hour.
	Timeout time.Duration
	// Callback function that should perform the authorization of the authenticated user. Called
	// only after an authentication success. Must return true on success, false on failure.
	// Optional, default to success.
	TokenValidator func(token *jwt.Token, c context.Context) error

	// Set the identity handler function. that mean the jwt is pass validete
	IdentityHandler func(c context.Context, claims jwt.MapClaims) (*Identity, error)

	// Private key file byte for asymmetric algorithms
	PrivKeyFile []byte

	// Public key file byte for asymmetric algorithms
	PubKeyFile []byte

	// Private key
	privKey *rsa.PrivateKey

	// Public key
	pubKey *rsa.PublicKey
}

var (
	// ErrMissingRealm indicates Realm name is required
	ErrMissingRealm = errors.New("realm is missing")

	// ErrForbidden when HTTP status 403 is given
	ErrForbidden = errors.New("you don't have permission to access this resource")

	// ErrInvalidPrivKey indicates that the given private key is invalid
	ErrInvalidClaims = errors.New("token payload content invalid")

	// ErrNoPrivKeyFile indicates that the given private key is unreadable
	ErrNoPrivKeyFile = errors.New("private key file unreadable")

	// ErrNoPubKeyFile indicates that the given public key is unreadable
	ErrNoPubKeyFile = errors.New("public key file unreadable")

	// ErrInvalidPrivKey indicates that the given private key is invalid
	ErrInvalidPrivKey = errors.New("private key invalid")

	// ErrInvalidPubKey indicates the the given public key is invalid
	ErrInvalidPubKey = errors.New("public key invalid")

	// ErrInvalidKey indicates the the given key is invalid
	ErrInvalidKey = errors.New("encrypty key invalid")
)

func (b *BearerTokenValidator) readKeys() error {
	if b.PrivKeyFile != nil {
		err := b.privateKey()
		if err != nil {
			return err
		}
	}
	if b.PubKeyFile != nil {
		err := b.publicKey()
		if err != nil {
			return err
		}
	}
	return nil
}

func (b *BearerTokenValidator) privateKey() error {
	if b.PrivKeyFile == nil {
		return ErrNoPrivKeyFile
	}
	key, err := jwt.ParseRSAPrivateKeyFromPEM(b.PrivKeyFile)
	if err != nil {
		return ErrInvalidPrivKey
	}
	b.privKey = key
	return nil
}

func (b *BearerTokenValidator) publicKey() error {
	if b.PubKeyFile == nil {
		return ErrNoPubKeyFile
	}
	key, err := jwt.ParseRSAPublicKeyFromPEM(b.PubKeyFile)
	if err != nil {
		return ErrInvalidPubKey
	}
	b.pubKey = key
	return nil
}

// Init initialize jwt configs.
func (b *BearerTokenValidator) Init() error {

	if err := b.readKeys(); err != nil {
		return err
	}

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
			if len(b.Key) == 0 {
				return nil, ErrInvalidKey
			}
			return b.Key, nil
		} else {
			if b.pubKey == nil {
				return nil, ErrInvalidPubKey
			}
			return b.pubKey, nil
		}
	})
}
