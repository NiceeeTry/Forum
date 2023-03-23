package models

import (
	"errors"
	"time"
)

var (
	ErrNoRecord           = errors.New("models: no matching record found")
	ErrInvalidCredentials = errors.New("models: invalid credentials")
	ErrDuplicateEmail     = errors.New("models: duplicate email")
	ErrDuplicateName      = errors.New("models: duplicate username")
)

type Post struct {
	ID       int
	Title    string
	Content  string
	Created  time.Time
	Name     string
	Category string
}

type User struct {
	ID             int
	Name           string
	Email          string
	HashedPassword []byte
	Created        time.Time
	Active         bool
}

type Comment struct {
	CommentText     string
	Author          string
	Created         time.Time
	LikesComment    int
	DislikesComment int
	ID              int
	PostId          int
}
