package persistance

import (
	"fmt"
	user "go-authentication/src/internal/core"
	"go-authentication/src/pkg"
)

type UserRepo interface {
	CreateUser(newUser user.User) (user.User, error)
	GetUser(username string) (user.User, error)
	MatchPassword(user user.User, password string) error
	CreateSession(session user.Session) error
}

type UserRepoPsql struct {
	db *Database
}

func NewUserRepo(d *Database) *UserRepoPsql {
	return &UserRepoPsql{db: d}
}

func (u *UserRepoPsql) CreateUser(newUser user.User) (user.User, error) {
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

func (u *UserRepoPsql) GetUser(username string) (user.User, error) {

	var newUser user.User
	query := "select uid, username, password from users where username = $1"
	err := u.db.db.QueryRow(query, username).Scan(&newUser.Uid, &newUser.Username, &newUser.Password)
	if err != nil {
		return user.User{}, err
	}
	return newUser, nil
}

func (u *UserRepoPsql) MatchPassword(user user.User, password string) error {

	err := pkg.CheckPassword(user.Password, password)

	if err != nil {
		fmt.Println(err, "unable to match password")
		return err
	}
	return nil
}

func (u *UserRepoPsql) CreateSession(session user.Session) error {
	_, err := u.db.db.Exec("INSERT INTO sessions (id, user_id, token_hash, expires_at, issued_at) VALUES ($1, $2, $3, $4, $5) ON CONFLICT (user_id) DO UPDATE SET token_hash = EXCLUDED.token_hash, expires_at = EXCLUDED.expires_at, issued_at = EXCLUDED.issued_at", session.Id, session.Uid, session.TokenHash, session.ExpiresAt, session.IssuedAt)
	if err != nil {
		return err
	}
	return nil
}
