package user

import "github.com/jmoiron/sqlx"

func IntializeHandler(pool *sqlx.DB) *UserHandler {
	userService := newUserService(newUserRepository(pool))
	return &UserHandler{
		signUpUsecase: userService,
		loginUsecase: userService,
	}
}
