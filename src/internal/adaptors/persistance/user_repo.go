package persistance

import (
	"fmt"
	user "go-authentication/src/internal/core/user"
	"go-authentication/src/pkg"
)
type UserRepo struct {
	db *Database
}

func NewUserRepo(d *Database) UserRepo {
	return UserRepo{db: d}
}

func (u *UserRepo) CreateUser(newUser user.User) (user.User, error) {
	var uid int
	query := "insert into users (username, password) values ($1, $2) returning uid"

	hpassword, err := pkg.HashPassword(newUser.Password)

	if err != nil {
		fmt.Println(err, "unable to hash password")
	}

	err = u.db.db.QueryRow(query, newUser.Username, hpassword).Scan(&uid)

	if err != nil {
		return user.User{}, err
	}

	newUser.Uid = uid

	return newUser, nil
}

func (u *UserRepo) GetUser(username string) (user.User, error) {
	var newUser user.User
	query := "select uid, username, password from users where username = $1"
	err := u.db.db.QueryRow(query, username).Scan(&newUser.Uid, &newUser.Username, &newUser.Password)
	if err != nil {
		return user.User{}, err
	}
	return newUser, nil
}

func (u *UserRepo) GetUserById(id int) (user.User, error) {
	var newUser user.User
	query := "select uid, username, password from users where uid = $1"
	err := u.db.db.QueryRow(query, id).Scan(&newUser.Uid, &newUser.Username, &newUser.Password)
	if err != nil {
		return user.User{}, err
	}
	return newUser, nil
}
