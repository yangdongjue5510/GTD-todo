package user_test

import (
	"testing"
	"time"

	"yangdongju/gtd_todo/internal/user"
	"yangdongju/gtd_todo/testhelper"

	"github.com/stretchr/testify/assert"
)

func TestFindUserByEmail(t *testing.T) {
	testDB := testhelper.GetTestDB()
	testhelper.CleanUp()

	email := "hello@example.com"
	password := "password1234_hash"
	_, _ = testDB.Exec("INSERT INTO users (email, password_hash) VALUES ($1, $2)", email, password)

	userRepository := user.NewUserRepository(testDB)

	foundUser, err := userRepository.FindUserByEmail(email)

	assert.Nil(t, err)
	assert.Equal(t, foundUser.ID, 1)
	assert.Equal(t, foundUser.Email, email)
}

func TestFindUserByEmail_NotFound(t *testing.T) {
	testDB := testhelper.GetTestDB()
	testhelper.CleanUp()
	email := "hello@example.com"
	userRepository := user.NewUserRepository(testDB)

	foundUser, err := userRepository.FindUserByEmail(email)

	assert.Nil(t, err)
	assert.Nil(t, foundUser)
}

func TestSave(t *testing.T) {
	testDB := testhelper.GetTestDB()
	testhelper.CleanUp()
	userRepository := user.NewUserRepository(testDB)

	createdUser := user.User{
		Email:        "hello@example.com",
		PasswordHash: "password_hash",
		CreatedAt:    time.Now(),
	}

	savedUser, _ := userRepository.Save(&createdUser)

	assert.Equal(t, 1, savedUser.ID)
	assert.Equal(t, savedUser.Email, createdUser.Email)
	assert.Equal(t, savedUser.PasswordHash, createdUser.PasswordHash)
}
