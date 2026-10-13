package models

var (
	RoleAdmin   = "admin"
	RoleStudent = "student"
	RoleTeacher = "teacher"
)

type User struct {
	Id       string `json:"id"`
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
	Id    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type Storage interface {
	AddUser(*User) error
	GetUserById(string) (*User, error)
	GetAllUsers() ([]User, error)
	DeleteUserById(string) error
	GetUserByEmail(string) (*User, error)
	Close()
	Ping() error
}

type PostgresReqs struct {
	User, Pass, Name, Host string
}
