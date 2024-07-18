package userservice

import (
	"fmt"
	"go-authentication/src/internal/adaptors/persistance"
	"go-authentication/src/internal/core/session"
	user "go-authentication/src/internal/core/user"
	"go-authentication/src/pkg"
	"time"

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
	// TODO: check if user is already registered
	newUser, err := u.userRepo.CreateUser(user)
	return newUser, err
}


type LoginResponse struct {
	FoundUser user.User
	TokenString string
	TokenExpire time.Time
	Session session.Session
}

func (u *UserService) LoginUser(requestUser user.User) (LoginResponse, error) {
	loginResponse := LoginResponse{}

	foundUser, err := u.userRepo.GetUser(requestUser.Username)
	if err != nil {
		return loginResponse, fmt.Errorf("invalid username")
	}
	
	loginResponse.FoundUser = foundUser
	if err := matchPassword(foundUser, requestUser.Password); err != nil {
		return loginResponse, fmt.Errorf("invalid password")
	}

	tokenString, tokenExpire, err := pkg.GenerateJWT(foundUser.Uid)
	loginResponse.TokenString = tokenString
	loginResponse.TokenExpire = tokenExpire

	if err != nil {
		return loginResponse, fmt.Errorf("failed to generate jwt")
	}

	session, err := pkg.GenerateSession(foundUser.Uid)
	loginResponse.Session = session
	if err != nil {
		return loginResponse, fmt.Errorf("failed to generate session")
	}

	err = u.sessionRepo.CreateSession(session)
	if err != nil {
		return loginResponse, fmt.Errorf("failed to create session")
	}
	
	return loginResponse, nil
}

func (u *UserService) GetJwtFromSession(sess string) (string, time.Time, error) {
	var tokenString string
	var tokenExpire time.Time
	session, err := u.sessionRepo.GetSession(sess)
	if err != nil {
		return tokenString, tokenExpire, err
	}

	err = matchSessionToken(sess, session.TokenHash)
	if err != nil {
		return tokenString, tokenExpire, err
	}

	tokenString, tokenExpire, err = pkg.GenerateJWT(session.Uid)
	if err != nil {
		return tokenString, tokenExpire, err
	}
	
	return tokenString, tokenExpire, nil
}

func (u *UserService) GetUserById(id int) (user.User, error) {
	newUser, err := u.userRepo.GetUserById(id)
	return newUser, err
}

func (u *UserService) LogoutUser(id int) error {
	err := u.sessionRepo.DeleteSession(id)
	return  err
}

func matchPassword(user user.User, password string) error {
	err := pkg.CheckPassword(user.Password, password)
	if err != nil {
		return fmt.Errorf("unable to match password: %w", err)
	}
	return nil
}

func matchSessionToken(id string, tokenHash string) error {
	err := bcrypt.CompareHashAndPassword([]byte(tokenHash), []byte(id))
	if err != nil {
		fmt.Println(err, "unable to match password")
		return err
	}
	return nil
}
