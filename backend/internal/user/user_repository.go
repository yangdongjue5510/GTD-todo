package user

import (
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
)

type UserRepository interface {
	FindUserByEmail(email string) (*User, error)
	Save(user *User) (*User, error)
}

type UserRepositoryImpl struct {
	db *sqlx.DB
}

type User struct {
	ID           int       `db:"id"`
	Email        string    `db:"email"`
	PasswordHash string    `db:"password_hash"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}

func NewUserRepository(db *sqlx.DB) *UserRepositoryImpl {
	return &UserRepositoryImpl{db: db}
}

func (r *UserRepositoryImpl) FindUserByEmail(email string) (*User, error) {
	var user User

	err := r.db.Get(&user, "SELECT * FROM users WHERE email = $1", email)
	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepositoryImpl) Save(user *User) (*User, error) {
	var id int
	err := r.db.QueryRow(`
		INSERT INTO users (email, password_hash, created_at)
		VALUES ($1, $2, $3)
		RETURNING id`,
		user.Email, user.PasswordHash, user.CreatedAt).Scan(&id)
	user.ID = id
	return user, err 
}
