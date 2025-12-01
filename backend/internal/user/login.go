package user

func (s *userService) login(request loginRequest) (*loginResponse, error) {
	return nil, nil
}

type loginUsecase interface {
	login(request loginRequest) (*loginResponse, error)
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginResponse struct {
	Token string `json:"token"`
}