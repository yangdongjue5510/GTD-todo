package user_test

import (
	"errors"
	"testing"

	"yangdongju/gtd_todo/internal/user"
	usermocks "yangdongju/gtd_todo/internal/user/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

// ============ Test Cases ============

func TestSignUp_Success(t *testing.T) {
	// given
	var capturedUser *user.User
	mockRepo := usermocks.NewUserRepository(t)
	mockIssuer := usermocks.NewIssuer(t)
	mockParser := usermocks.NewParser(t)

	request := user.SignUpRequest{
		Email:    "newuser@example.com",
		Password: "password1234",
	}

	// Mock 설정
	mockRepo.EXPECT().FindUserByEmail(request.Email).Return(nil, nil)
	mockRepo.EXPECT().Save(mock.MatchedBy(func(u *user.User) bool {
		capturedUser = u
		return u.Email == request.Email &&
			bcrypt.CompareHashAndPassword([]byte(u.PasswordHash),
				[]byte(request.Password)) == nil
	})).RunAndReturn(func(u *user.User) (*user.User, error) {
		return &user.User{
			ID:           1,
			Email:        u.Email,
			PasswordHash: u.PasswordHash,
			CreatedAt:    u.CreatedAt,
		}, nil
	})

	service := user.NewUserService(mockRepo, mockIssuer, mockParser)

	// when
	response, err := service.SignUp(request)

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
	mockRepo := usermocks.NewUserRepository(t)
	mockIssuer := usermocks.NewIssuer(t)
	mockParser := usermocks.NewParser(t)

	existingUser := &user.User{
		ID:    10,
		Email: "existing@example.com",
	}

	request := user.SignUpRequest{
		Email:    "existing@example.com",
		Password: "password1234",
	}

	// Mock 설정
	mockRepo.EXPECT().FindUserByEmail(request.Email).Return(existingUser, nil)

	service := user.NewUserService(mockRepo, mockIssuer, mockParser)

	// when
	response, err := service.SignUp(request)

	// then
	assert.Nil(t, response)
	assert.NotNil(t, err)

	var userAlreadyExistsErr *user.UserAlreadyExistsError
	assert.True(t, errors.As(err, &userAlreadyExistsErr), "Error should be userAlreadyExistsError type")
	assert.Contains(t, err.Error(), "existing@example.com")
	assert.Contains(t, err.Error(), "10")
}

func TestSignUp_FindUserByEmail_RepositoryError(t *testing.T) {
	// given
	mockRepo := usermocks.NewUserRepository(t)
	mockIssuer := usermocks.NewIssuer(t)
	mockParser := usermocks.NewParser(t)

	repositoryError := errors.New("database connection failed")

	request := user.SignUpRequest{
		Email:    "test@example.com",
		Password: "password1234",
	}

	// Mock 설정
	mockRepo.EXPECT().FindUserByEmail(request.Email).Return(nil, repositoryError)

	service := user.NewUserService(mockRepo, mockIssuer, mockParser)

	// when
	response, err := service.SignUp(request)

	// then
	assert.Nil(t, response)
	assert.NotNil(t, err)
	assert.Equal(t, "database connection failed", err.Error())
}

func TestSignUp_Save_RepositoryError(t *testing.T) {
	// given
	mockRepo := usermocks.NewUserRepository(t)
	mockIssuer := usermocks.NewIssuer(t)
	mockParser := usermocks.NewParser(t)

	saveError := errors.New("failed to insert user")

	request := user.SignUpRequest{
		Email:    "test@example.com",
		Password: "password1234",
	}

	// Mock 설정
	mockRepo.EXPECT().FindUserByEmail(request.Email).Return(nil, nil)
	mockRepo.EXPECT().Save(mock.Anything).Return(nil, saveError)

	service := user.NewUserService(mockRepo, mockIssuer, mockParser)

	// when
	response, err := service.SignUp(request)

	// then
	assert.Nil(t, response)
	assert.NotNil(t, err)
	assert.Equal(t, "failed to insert user", err.Error())
}
