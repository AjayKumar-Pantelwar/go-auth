package userservice

import (
	"go-authentication/src/internal/adaptors/persistance"
)

type UserService struct {
	persistance.UserRepo
}


func NewUserService(userRepo persistance.UserRepo) UserService {
	return UserService{userRepo}
}

type UserServiceImpl interface {
	persistance.UserRepo
}
