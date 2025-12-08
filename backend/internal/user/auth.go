package user

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID int    `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

type Issuer interface {
	Issue(userId int, email string, duration time.Duration) (string, error)
}

type Parser interface {
	Parse(token string) (*Claims, error)
}

type tokenService struct {
	secretKey []byte
	now       func() time.Time
}

func NewTokenService(secretKey string, now func() time.Time) (*tokenService, error) {
	if secretKey == "" {
		return nil, errors.New("secret key should not be empty")
	}

	if now == nil {
		return nil, errors.New("now function shoud not be nil")
	}
	return &tokenService{
		secretKey: []byte(secretKey),
		now:       now,
	}, nil
}

func (service *tokenService) Issue(userID int, email string, duration time.Duration) (string, error) {
	now := service.now()
	claims := &Claims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(duration)),
			IssuedAt:  jwt.NewNumericDate(now),
			Issuer:    "gtd-todo-app",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(service.secretKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}

func (service *tokenService) Parse(tokenString string) (*Claims, error) {
	parser := jwt.NewParser(
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
		jwt.WithIssuer("gtd-todo-app"),
		jwt.WithTimeFunc(service.now), // IssuedAt/ExpiresAt 검증 시 기준 시간
	)
	token, err := parser.ParseWithClaims(tokenString, &Claims{}, func(t *jwt.Token) (any, error) {
		return service.secretKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, errors.New("failed to cast token claims")
	}

	if !token.Valid {
		return nil, jwt.ErrSignatureInvalid
	}

	return claims, nil
}