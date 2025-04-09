package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTClaims struct {
	//(issuer): issuer of the JWT
	iss string
	//(subject): Subject of the JWT (the user)
	sub string
	// (audience): Recipient for which the JWT is intended
	aud string
	// (expiration time): Time after which the JWT expires
	exp time.Time
	// (not before time): Time before which the JWT must not be accepted for processing
	nbf time.Time
	// (issued at time): Time at which the JWT was issued; can be used to determine age of the JWT
	iat time.Time
	// jti (JWT ID): Unique identifier; can be used to prevent the JWT from being replayed (allows a token to be used only once)
	// jti
}

type JWTAuthenticator struct {
	jwtClaims JWTClaims
	secret    string
}

// simplified JWTAuthenticator constructor with secret, aud = sub, iss, exp and nbf=iat
func NewSimpleJWTAuthenticator(secret string, iss, aud string) *JWTAuthenticator {

	return &JWTAuthenticator{
		JWTClaims{
			iss: iss,
			aud: aud,
		},
		secret,
	}
}

func (auth *JWTAuthenticator) GenerateToken(claims jwt.Claims) (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(auth.secret))

	if err != nil {
		return "", err
	}

	return tokenString, err
}

func (auth *JWTAuthenticator) ValidateToken(token string) (*jwt.Token, error) {

	return jwt.Parse(
		token,
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method %v", t.Header["alg"])
			}
			return []byte(auth.secret), nil
		},
		// options
		jwt.WithExpirationRequired(),
		jwt.WithAudience(auth.jwtClaims.aud),
		jwt.WithIssuer(auth.jwtClaims.aud),
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Name}),
	)
}
