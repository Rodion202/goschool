package models

import "context"

const (
	RoleAdmin   = "admin"
	RoleStudent = "student"
	RoleTeacher = "teacher"
)

type Config struct {
	PGXuser        string
	PGXpass        string
	PGXdb          string
	JWTsecret      string
	PGXHost        string
	ADMIN_EMAIL    string
	ADMIN_NAME     string
	ADMIN_PASSWORD string
}

type User struct {
	Id       int64  `json:"id"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type LoginInput struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}
type UserMessage struct {
	Id    int64  `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

type Storage interface {
	AddUser(context.Context, *User) error
	GetUserById(context.Context, int64) (*User, error)
	GetAllUsers(context.Context) ([]User, error)
	DeleteUserById(context.Context, int64) error
	GetUserByEmail(context.Context, string) (*User, error)
}
