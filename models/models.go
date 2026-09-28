package models

type User struct{
	Id string `json:"id"`
	Email string `json:"email"`
	Name string `json:"name"`
	Password string `json:"-"`
}

type UserMessage struct{
	Id string
	Email string
}

type Storage interface{
 AddUser(User)error
 GetUserById(string)(*User,error)
 GetAllUsers()([]User, error)
 DeleteUserById(string)error
 Close()
 Ping()error
}

type PostgresReqs struct{
	User,Pass,Name,Host string
}
