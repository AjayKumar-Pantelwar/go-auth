package userservice

import (
	"fmt"
	"go-authentication/src/internal/adaptors/persistance"
	session "go-authentication/src/internal/core/session"
	user "go-authentication/src/internal/core/user"
	"go-authentication/src/pkg"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepo persistance.UserRepo
	sessionRepo persistance.SessionRepo
}


func NewUserService(userRepo persistance.UserRepo, sessionRepo persistance.SessionRepo) UserService {
	return UserService{userRepo: userRepo, sessionRepo: sessionRepo}
}

func (u *UserService) RegisterUser(user user.User) (user.User, error) {
	newUser, err := u.userRepo.CreateUser(user)
	return newUser, err
}

func (u *UserService) GetUser(username string) (user.User, error) {
	newUser, err := u.userRepo.GetUser(username)
	return newUser, err
}

func (u *UserService) MatchPassword(user user.User, password string) error {
	err := pkg.CheckPassword(user.Password, password)
	if err != nil {
		return fmt.Errorf("unable to match password: %w", err)
	}
	return nil
}

func (u *UserService) CreateSession(session session.Session) error {
	err := u.sessionRepo.CreateSession(session)
	return err
}

func (u *UserService) GetUserById(id int) (user.User, error) {
	newUser, err := u.userRepo.GetUserById(id)
	return newUser, err
}

func (u *UserService) GetSession(id string) (session.Session, error) {
	session, err := u.sessionRepo.GetSession(id)
	return session, err
}

func (u *UserService) MatchSessionToken(id string, tokenHash string) error {
	err := bcrypt.CompareHashAndPassword([]byte(tokenHash), []byte(id))
	if err != nil {
		fmt.Println(err, "unable to match password")
		return err
	}
	return nil
}

func (u *UserService) DeleteSession(id int) error {
	err := u.sessionRepo.DeleteSession(id)
	return  err
}