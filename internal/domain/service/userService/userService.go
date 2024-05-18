package userservice

type UserService struct {
}

func New() *UserService {
	return &UserService{}
}

func (us *UserService) Login() {}

func (us *UserService) Registration() {}

func (us *UserService) GetUser() {}

func (us *UserService) GetAllUsers() {}

func (us *UserService) UpdateUser() {}
