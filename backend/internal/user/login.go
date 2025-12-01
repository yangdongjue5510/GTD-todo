package user

func (s *userService) login(request loginRequest) (*loginResponse, error) {
	return nil, nil
}

type loginUsecase interface {
	login(request loginRequest) (*loginResponse, error)
}

type loginRequest struct {
	email    string `json:email`
	password string `json:password`
}

type loginResponse struct {
	token string `json:token`
}