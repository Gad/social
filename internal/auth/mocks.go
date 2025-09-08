package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type MockAuthenticator struct {
	jwtClaims JWTClaims
	secret    string
}

var testSecret = []byte("testsecret")

var testClaims = jwt.MapClaims{
	"iss": "test-issuer",
	"aud": "test-audience",
	"sub": int64(42),
	"exp": time.Now().Add(time.Hour).Unix(),
	"nbf": time.Now().Unix(),
	"iat": time.Now().Unix(),
}

func (m MockAuthenticator) GenerateToken(claims jwt.Claims) (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, testClaims)

	tokenString, err := token.SignedString(testSecret)

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (m MockAuthenticator) ValidateToken(token string) (*jwt.Token, error) {
	tokenObj, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return testSecret, nil
	})
	if err != nil {
		return nil, ErrInvalidToken
	}
	if !tokenObj.Valid {
		return nil, ErrInvalidToken
	}
	return tokenObj, nil
}
