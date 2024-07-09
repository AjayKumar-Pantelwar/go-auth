package user

import "go-authentication/src/internal/adaptors/persistance"

type UserService struct {
	userRepo persistance.UserRepo
}

func NewUserService(userRepo persistance.UserRepo) UserService {
	return UserService{userRepo: userRepo}
}
