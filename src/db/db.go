package db

import (
	"database/sql"
	"fmt"
	util "go-authentication/src/utils"

	_ "github.com/lib/pq"
)

type Database struct {
	db *sql.DB
}

func NewDatabase() (*Database, error) {
	db, err := sql.Open("postgres", "postgresql://postgres:postgres@localhost:5432/goauth?sslmode=disable")
	if err != nil {
		return nil, err
	}
	fmt.Println(("Connected to database"))
	return &Database{db: db}, nil
}

func (d *Database) Close() {
	d.db.Close()
}

func (d *Database) GetDB() *sql.DB {
	return d.db
}

func CreateUser(d *Database, user User) (User, error) {

	var uid int
	query := "insert into users (username, password) values ($1, $2) returning uid"

	hpassword, err := util.HashPassword(user.Password)

	if err != nil {
		fmt.Println(err, "unable to hash password")
	}

	err = d.db.QueryRow(query, user.Username, hpassword).Scan(&uid)

	if err != nil {
		return User{}, err
	}

	user.Uid = uid

	return user, nil

}

func GetUser(d *Database, username string) (User, error) {

	var user User
	query := "select uid, username, password from users where username = $1"
	err := d.db.QueryRow(query, username).Scan(&user.Uid, &user.Username, &user.Password)
	if err != nil {
		return User{}, err
	}
	return user, nil
}

func MatchPassword(d *Database, user User, password string) error {

	err := util.CheckPassword(user.Password, password)

	if err != nil {
		fmt.Println(err, "unable to match password")
		return err
	}
	return nil
}
