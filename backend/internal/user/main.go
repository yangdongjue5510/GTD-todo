//go:generate mockery
package user

import (
	"log"
	"os"
	"time"

	"github.com/jmoiron/sqlx"
)

func IntializeHandler(pool *sqlx.DB) *UserHandler {
	tokenService := initTokenService()
	userService := NewUserService(NewUserRepository(pool), tokenService, tokenService)
	return &UserHandler{
		signUpUsecase: userService,
		loginUsecase:  userService,
	}
}

func initTokenService() *tokenService {
	JWTSecretKey := os.Getenv("JWT_SECRET_KEY")
	timeFunc := func() time.Time { return time.Now() }
	tokenService, err := NewTokenService(JWTSecretKey, timeFunc)
	if err != nil {
		log.Fatalf("TokenService init failed. %v\n", err)
	}
	return tokenService
}
