package main

import (
	"fmt"
	"os"

	"goschool/database"
	"goschool/handlers"

	"github.com/gin-gonic/gin"
)

func hello(c *gin.Context) {
	c.JSON(200, gin.H{"message": "hello!"})
}

func hasDropFlag() bool {
	if len(os.Args) < 2 {
		return false
	}
	return os.Args[1] == "drop"
}

func main() {
	storage, err := database.Connect(hasDropFlag())
	if err != nil {
		fmt.Println("Error conecting Postgres:", err)
	}
	defer storage.Close()

	hs := handlers.NewStorageHandler(storage)

	r := gin.Default()

	r.GET("/", hello)

	usersRouter := r.Group("/users")
	{
		usersRouter.GET("/", hs.GetAllUsers)
		usersRouter.POST("/", hs.AddUser)
		usersRouter.POST("/login", hs.Login)
	}

	fmt.Println("http://localhost:8080/")
	if err := r.Run(":8080"); err != nil {
		fmt.Println(err)
		panic(err)
	}
}
