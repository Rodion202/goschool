package models

import "context"

var (
	RoleAdmin   = "admin"
	RoleStudent = "student"
	RoleTeacher = "teacher"
)

type User struct {
	Id       int64  `json:"id"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type UserMessage struct {
	Id    int64  `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

type Storage interface {
	AddUser(context.Context, *User) error
	GetUserById(context.Context, string) (*User, error)
	GetAllUsers(context.Context) ([]User, error)
	DeleteUserById(context.Context, string) error
	GetUserByEmail(context.Context, string) (*User, error)
	Close()
	Ping() error
}

type PostgresReqs struct {
	User, Pass, Name, Host string
}
