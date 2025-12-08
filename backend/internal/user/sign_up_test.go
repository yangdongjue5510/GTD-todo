package user

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

// ============ Mock UserRepository ============
type mockUserRepository struct {
	findUserByEmailStub func(email string) (*User, error)
	saveStub            func(user *User) (*User, error)
}

func (m *mockUserRepository) FindUserByEmail(email string) (*User, error) {
	if m.findUserByEmailStub != nil {
		return m.findUserByEmailStub(email)
	}
	return nil, errors.New("findUserByEmailStub not implemented")
}

func (m *mockUserRepository) Save(user *User) (*User, error) {
	if m.saveStub != nil {
		return m.saveStub(user)
	}
	return nil, errors.New("saveStub not implemented")
}

// ============ Test Cases ============

func TestSignUp_Success(t *testing.T) {
	// given
	var capturedUser *User
	mockRepo := &mockUserRepository{
		findUserByEmailStub: func(email string) (*User, error) {
			return nil, nil
		},
		saveStub: func(user *User) (*User, error) {
			capturedUser = user
			user.ID = 1
			return user, nil
		},
	}

	service := &userService{userRepository: mockRepo}
	request := SignUpRequest{
		Email:    "newuser@example.com",
		Password: "password1234",
	}

	// when
	response, err := service.signUp(request)

	// then
	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, 1, response.ID)
	assert.Equal(t, "newuser@example.com", response.Email)

	// 패스워드 해싱 검증
	assert.NotNil(t, capturedUser)
	assert.NotEqual(t, "password1234", capturedUser.PasswordHash, "Password should be hashed")
	err = bcrypt.CompareHashAndPassword([]byte(capturedUser.PasswordHash), []byte("password1234"))
	assert.Nil(t, err, "Password hash should match original password")

	// CreatedAt 검증
	assert.NotNil(t, capturedUser.CreatedAt)
}

func TestSignUp_UserAlreadyExists(t *testing.T) {
	// given
	mockRepo := &mockUserRepository{
		findUserByEmailStub: func(email string) (*User, error) {
			return &User{
				ID:    10,
				Email: email,
			}, nil
		},
	}

	service := &userService{userRepository: mockRepo}
	request := SignUpRequest{
		Email:    "existing@example.com",
		Password: "password1234",
	}

	// when
	response, err := service.signUp(request)

	// then
	assert.Nil(t, response)
	assert.NotNil(t, err)

	var userAlreadyExistsErr *userAlreadyExistsError
	assert.True(t, errors.As(err, &userAlreadyExistsErr), "Error should be userAlreadyExistsError type")
	assert.Contains(t, err.Error(), "existing@example.com")
	assert.Contains(t, err.Error(), "10")
}

func TestSignUp_FindUserByEmail_RepositoryError(t *testing.T) {
	// given
	mockRepo := &mockUserRepository{
		findUserByEmailStub: func(email string) (*User, error) {
			return nil, errors.New("database connection failed")
		},
	}

	service := &userService{userRepository: mockRepo}
	request := SignUpRequest{
		Email:    "test@example.com",
		Password: "password1234",
	}

	// when
	response, err := service.signUp(request)

	// then
	assert.Nil(t, response)
	assert.NotNil(t, err)
	assert.Equal(t, "database connection failed", err.Error())
}

func TestSignUp_Save_RepositoryError(t *testing.T) {
	// given
	mockRepo := &mockUserRepository{
		findUserByEmailStub: func(email string) (*User, error) {
			return nil, nil
		},
		saveStub: func(user *User) (*User, error) {
			return nil, errors.New("failed to insert user")
		},
	}

	service := &userService{userRepository: mockRepo}
	request := SignUpRequest{
		Email:    "test@example.com",
		Password: "password1234",
	}

	// when
	response, err := service.signUp(request)

	// then
	assert.Nil(t, response)
	assert.NotNil(t, err)
	assert.Equal(t, "failed to insert user", err.Error())
}
