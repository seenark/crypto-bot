package usersService

import (
	"fmt"
	"mybot/repository/users"
)

type usersService struct {
	repo users.UserRepository
}

func NewUserService(repo users.UserRepository) UserService {
	return &usersService{
		repo: repo,
	}
}

func (usr usersService) Create(u users.User) (string, error) {
	user, err := usr.repo.Create(u)
	if err != nil {
		return "", fmt.Errorf("cannot create new user %v", u)
	}
	return user, nil
}
func (usr usersService) GetAll() ([]UserResponse, error) {
	users, err := usr.repo.GetAll()
	if err != nil {
		return nil, fmt.Errorf("error while fetching all user")
	}
	userRes := []UserResponse{}
	for _, user := range users {
		userRes = append(userRes, UserResponse(user))
	}
	return userRes, err
}
func (usr usersService) GetUserById(id string) (*UserResponse, error) {
	user, err := usr.repo.GetUserById(id)
	if err != nil {
		return nil, err
	}
	newUser := UserResponse(*user)
	return &newUser, nil
}
func (usr usersService) UpdateUserById(id string, u users.User) (*UserResponse, error) {
	updated, err := usr.repo.UpdateUserById(id, u)
	if err != nil {
		return nil, err
	}
	newUser := UserResponse(*updated)
	return &newUser, nil
}
func (usr usersService) DeleteUserById(id string) (int, error) {
	deleted, err := usr.repo.DeleteUserById(id)
	if err != nil {
		return 0, err
	}
	return deleted, nil
}
