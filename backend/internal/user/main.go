package user

import "github.com/jmoiron/sqlx"

func IntializeHandler(pool *sqlx.DB) *UserHandler {
	return newUserHandler(newUserService(newUserRepository(pool)))
}
