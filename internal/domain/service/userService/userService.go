package userservice

type UserService struct {
}

func New() *UserService {
	return &UserService{}
}

func (us *UserService) Login() {}

func (us *UserService) Registration() {}
