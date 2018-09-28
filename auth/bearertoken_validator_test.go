//
// 可以去https://jwt.io生成
//
package auth_test

import (
	"context"
	"github.com/dgrijalva/jwt-go"
	"github.com/qeelyn/go-common/auth"
	"testing"
)

func TestBearTokenValidator_ValidateHS(t *testing.T) {
	username, _ := 18819, "sun5kong"
	HSexpToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1MzgxMDIwODAsImlhdCI6MTUzODEwMDI4MCwic3ViIjoiMTg4MTkifQ.hxBaH6S4cjYI4IGTbVzryGTJvvZ0xJ_KA5vL_i1ekJA"
	key := "passw0rd"
	HStokenNoexp := "eyJhbGciOiJIUzI1NiJ9.eyJpYXQiOjE1MzgwNTE4ODAsInN1YiI6IjE4ODE5In0.vV7V-B2TZQnu9s24YWnq4IKjC1OJfkBnyzBpAJ1kVDg"
	b := &auth.BearerTokenValidator{
		Key: []byte(key),
		IdentityHandler: func(c context.Context, claims jwt.MapClaims) (*auth.Identity, error) {
			id := claims["sub"].(string)
			return &auth.Identity{
				Id: id,
			}, nil
		},
	}
	id, err := b.Validate(context.TODO(), HStokenNoexp)
	if err != nil {
		t.Fatal(err)
	}
	if id.IdInt() != int32(username) {
		t.Error("id wrong")
	}

	_, err = b.Validate(context.TODO(), HSexpToken)
	if err == nil {
		t.Fatal(err)
	}
}

func TestBearTokenValidator_ValidateRS(t *testing.T) {
	username := 122
	RSToken := "eyJhbGciOiJSUzI1NiJ9.eyJleHAiOjE1MzkwNTkwODAsImlhdCI6MTUzODA1MTg4MCwic3ViIjoiMTIyIn0.ogNGcmT2I6aET15cXrKbG3HUdG-mOgvYeJh8P-vxu2lk9oYyIqnouYv_PR2oKsEasEBjcKpDTS7NLh4GMkUx-NT0WyqbPIw3ZAcnsWruJg4fpk0bf75JXcl5gs1STCHaNlh0JFsZNUJynTBgiyweeLqFlxQIpQaMr5DEGbA5yEY"
	RSTokenExp := "eyJhbGciOiJSUzI1NiJ9.eyJleHAiOjE1MzgwNTI4ODAsImlhdCI6MTUzODA1MTg4MCwic3ViIjoiMTIyIn0.np1h02qA7sVZbumXWGGPI8l7pkDCeFWuVwhd0YCyLPYSCfZ2DgIRikwEpLdPXrxcUCGJ0NYCzZpVGAMEESndJ7m0d8CbU6mICRuPpMYD3JNSR-qNbV-uiNoOFnk9sTFI_6t79W_0V9N1AlefzFwBtyfQAYAW_A0ByyWMKXMLgKg"
	key := "passw0rd"
	b := &auth.BearerTokenValidator{
		Key: []byte(key),
		PubKeyFile: []byte(
			`-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDutm+bsYoI7vy2PSvQQEovcVNp
9sAkZAxUbSpKsE/VGNxhEwHP/LaTLK8mX7Yo1FeA6kWxkD/s0w07YUeiIZc4D2Hd
UjrhTZm6Fn6ZyziuMoinEG82Rz0B0ggqDMiGj73SQBdGRqFUALL0EwStRNwBEta3
3+Od61UJ1mqf7ZArhwIDAQAB
-----END PUBLIC KEY-----
`),
		IdentityHandler: func(c context.Context, claims jwt.MapClaims) (*auth.Identity, error) {
			id := claims["sub"].(string)
			return &auth.Identity{
				Id: id,
			}, nil
		},
	}
	err := b.Init()
	if err != nil {
		t.Fatal(err)
	}
	id, err := b.Validate(context.TODO(), RSToken)
	if err != nil {
		t.Fatal(err)
	}
	if id.IdInt() != int32(username) {
		t.Error("id wrong")
	}

	_, err = b.Validate(context.TODO(), RSTokenExp)
	if err == nil {
		t.Fatal("toekn exp error")
	}
}

func TestBearTokenValidator_ValidateComplex(t *testing.T) {
	username := 122
	key := "passw0rd"
	RStoken := "eyJhbGciOiJSUzI1NiJ9.eyJpYXQiOjE1MzgwNTE4ODAsInN1YiI6IjEyMiJ9.MRMOEvVuBO3h0ZuQJ9Pxt4CmJN8_eUWDZTmAkZemmCzZixVNn1l03s2Hjs7NuoIm1KRSACr3WJvN4kozKhydjXYhalU5cj2HBfEopE1NKkqsbkjdAfNGwJy-kxKGn7edVc6DQ1Nx_3gcqMcDPeef5DWmHJR0-dhKV65aqT6O7LQ"
	HStoken := "eyJhbGciOiJIUzI1NiJ9.eyJpYXQiOjE1MzgwNTE4ODAsInN1YiI6IjEyMiJ9.zxOF3-wx7xEwWXK2H3aWSKOKups6XZZO-rpYSN985AA"
	b := &auth.BearerTokenValidator{
		Key: []byte(key),
		PubKeyFile: []byte(
			`-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDutm+bsYoI7vy2PSvQQEovcVNp
9sAkZAxUbSpKsE/VGNxhEwHP/LaTLK8mX7Yo1FeA6kWxkD/s0w07YUeiIZc4D2Hd
UjrhTZm6Fn6ZyziuMoinEG82Rz0B0ggqDMiGj73SQBdGRqFUALL0EwStRNwBEta3
3+Od61UJ1mqf7ZArhwIDAQAB
-----END PUBLIC KEY-----
`),
		IdentityHandler: func(c context.Context, claims jwt.MapClaims) (*auth.Identity, error) {
			id := claims["sub"].(string)
			return &auth.Identity{
				Id: id,
			}, nil
		},
	}
	err := b.Init()
	if err != nil {
		t.Fatal(err)
	}
	// hs
	if _, err := b.Validate(context.TODO(), HStoken); err != nil {
		t.Error(err)
	}
	// rs
	id, err := b.Validate(context.TODO(), RStoken)
	if err != nil {
		t.Fatal(err)
	}
	if id.IdInt() != int32(username) {
		t.Error("id wrong")
	}
	// wrong rs
	RStoken = "it worng"
	if _, err := b.Validate(context.TODO(), RStoken); err == nil {
		t.Fatal("wrong token pass")
	}
}
