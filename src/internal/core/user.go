package user

import "time"

type User struct {
	Uid      int
	Username string
	Password string
}

type Session struct {
	Id        int
	Uid       int
	TokenHash string
	ExpiresAt time.Time
	IssuedAt  time.Time
}
