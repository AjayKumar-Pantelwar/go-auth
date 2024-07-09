package persistance

import (
	"fmt"
	user "go-authentication/src/internal/core"
	util "go-authentication/src/utils"
)

type UserRepo interface {
	CreateUser(u *UserRepoPsql, newUser user.User) (user.User, error)
	GetUser(u *UserRepoPsql, username string) (user.User, error)
	MatchPassword(u *UserRepoPsql, user user.User, password string) error
	CreateSession(u *UserRepoPsql, session user.Session) error
}

type UserRepoPsql struct {
	db *Database
}

func NewUserRepo(d *Database) *UserRepoPsql {
	return &UserRepoPsql{db: d}
}

func (*UserRepoPsql) CreateUser(u *UserRepoPsql, newUser user.User) (user.User, error) {
	var uid int
	query := "insert into users (username, password) values ($1, $2) returning uid"

	hpassword, err := util.HashPassword(newUser.Password)

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

func (*UserRepoPsql) GetUser(u *UserRepoPsql, username string) (user.User, error) {

	var newUser user.User
	query := "select uid, username, password from users where username = $1"
	err := u.db.db.QueryRow(query, username).Scan(&newUser.Uid, &newUser.Username, &newUser.Password)
	if err != nil {
		return user.User{}, err
	}
	return newUser, nil
}

func (*UserRepoPsql) MatchPassword(u *UserRepoPsql, user user.User, password string) error {

	err := util.CheckPassword(user.Password, password)

	if err != nil {
		fmt.Println(err, "unable to match password")
		return err
	}
	return nil
}

func (*UserRepoPsql) CreateSession(u *UserRepoPsql, session user.Session) error {
	_, err := u.db.db.Exec("INSERT INTO sessions (id, user_id, token_hash, expires_at, issued_at) VALUES ($1, $2, $3, $4, $5) ON CONFLICT (user_id) DO UPDATE SET token_hash = EXCLUDED.token_hash, expires_at = EXCLUDED.expires_at, issued_at = EXCLUDED.issued_at", session.Id, session.Uid, session.TokenHash, session.ExpiresAt, session.IssuedAt)
	if err != nil {
		return err
	}
	return nil
}
