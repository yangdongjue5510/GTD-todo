package user_test

import (
	"errors"
	"testing"
	"time"

	"yangdongju/gtd_todo/internal/user"
	usermocks "yangdongju/gtd_todo/internal/user/mocks"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

// ============ Login Service Tests ============

func TestLogin_Success(t *testing.T) {
	// given

	email := "test@example.com"
	password := "password123"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	mockRepo := usermocks.NewUserRepository(t)
	mockIssuer := usermocks.NewIssuer(t)
	mockParser := usermocks.NewParser(t)

	expectedUser := &user.User{
		ID:           1,
		Email:        email,
		PasswordHash: string(hashedPassword),
	}
	expectedToken := "generated.jwt.token"

	// Mock 설정
	mockRepo.EXPECT().FindUserByEmail(email).Return(expectedUser, nil)
	mockIssuer.EXPECT().Issue(expectedUser.ID, expectedUser.Email, 24*time.Hour).Return(expectedToken, nil)

	service := user.NewUserService(mockRepo, mockIssuer, mockParser)
	request := user.LoginRequest{
		Email:    email,
		Password: password,
	}

	// when
	response, err := service.Login(request)

	// then
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, expectedToken, response.Token)
}

func TestLogin_UserNotFound(t *testing.T) {
	// given
	email := "notfound@example.com"
	password := "password123"

	mockRepo := usermocks.NewUserRepository(t)
	mockIssuer := usermocks.NewIssuer(t)
	mockParser := usermocks.NewParser(t)

	// Mock 설정 - 사용자를 찾지 못함
	mockRepo.EXPECT().FindUserByEmail(email).Return(nil, nil)

	service := user.NewUserService(mockRepo, mockIssuer, mockParser)
	request := user.LoginRequest{
		Email:    email,
		Password: password,
	}

	// when
	response, err := service.Login(request)

	// then
	assert.Error(t, err)
	assert.Nil(t, response)

	var invalidCredError *user.InvalidCredentialsError
	assert.ErrorAs(t, err, &invalidCredError)
	assert.Equal(t, "Invalid email or password", err.Error())
}

func TestLogin_RepositoryError(t *testing.T) {
	// given
	email := "test@example.com"
	password := "password123"
	repositoryError := errors.New("database connection failed")

	mockRepo := usermocks.NewUserRepository(t)
	mockIssuer := usermocks.NewIssuer(t)
	mockParser := usermocks.NewParser(t)

	// Mock 설정 - Repository에서 에러 발생
	mockRepo.EXPECT().FindUserByEmail(email).Return(nil, repositoryError)

	service := user.NewUserService(mockRepo, mockIssuer, mockParser)
	request := user.LoginRequest{
		Email:    email,
		Password: password,
	}

	// when
	response, err := service.Login(request)

	// then
	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Equal(t, repositoryError, err)
}

func TestLogin_InvalidPassword(t *testing.T) {
	// given
	email := "test@example.com"
	correctPassword := "password123"
	wrongPassword := "wrongpassword"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(correctPassword), bcrypt.DefaultCost)

	mockRepo := usermocks.NewUserRepository(t)
	mockIssuer := usermocks.NewIssuer(t)
	mockParser := usermocks.NewParser(t)

	expectedUser := &user.User{
		ID:           1,
		Email:        email,
		PasswordHash: string(hashedPassword),
	}

	// Mock 설정
	mockRepo.EXPECT().FindUserByEmail(email).Return(expectedUser, nil)

	service := user.NewUserService(mockRepo, mockIssuer, mockParser)
	request := user.LoginRequest{
		Email:    email,
		Password: wrongPassword,
	}

	// when
	response, err := service.Login(request)

	// then
	assert.Error(t, err)
	assert.Nil(t, response)

	var invalidCredError *user.InvalidCredentialsError
	assert.ErrorAs(t, err, &invalidCredError)
	assert.Equal(t, "Invalid email or password", err.Error())
}

func TestLogin_TokenIssueError(t *testing.T) {
	// given
	email := "test@example.com"
	password := "password123"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	tokenError := errors.New("failed to sign token")

	mockRepo := usermocks.NewUserRepository(t)
	mockIssuer := usermocks.NewIssuer(t)
	mockParser := usermocks.NewParser(t)

	expectedUser := &user.User{
		ID:           1,
		Email:        email,
		PasswordHash: string(hashedPassword),
	}

	// Mock 설정
	mockRepo.EXPECT().FindUserByEmail(email).Return(expectedUser, nil)
	mockIssuer.EXPECT().Issue(expectedUser.ID, expectedUser.Email, 24*time.Hour).Return("", tokenError)

	service := user.NewUserService(mockRepo, mockIssuer, mockParser)
	request := user.LoginRequest{
		Email:    email,
		Password: password,
	}

	// when
	response, err := service.Login(request)

	// then
	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Equal(t, tokenError, err)
}
