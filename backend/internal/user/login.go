package user

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

func (s *userService) Login(request LoginRequest) (*LoginResponse, error) {
	user, err := s.userRepository.FindUserByEmail(request.Email)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, NewInvalidCredentialsError()
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(request.Password),
	)
	if err != nil {
		return nil, NewInvalidCredentialsError()
	}

	// 3. JWT 토큰 생성
	token, err := s.tokenIssuer.Issue(user.ID, user.Email, 24*time.Hour)
	if err != nil {
		return nil, err
	}

	return &LoginResponse{
		Token: token,
	}, nil
}

type LoginUsecase interface {
	Login(request LoginRequest) (*LoginResponse, error)
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}
