package main

import (
	"context"
	"fmt"
	"goschool/internal/database"
	"goschool/internal/envget"
	handlers "goschool/internal/handlers"
	"goschool/internal/users"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	if err := Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func hello(c *gin.Context) {
	c.JSON(200, gin.H{"message": "hello!"})
}

func hasDropFlag() bool {
	if len(os.Args) < 2 {
		return false
	}
	return os.Args[1] == "drop"
}

func Run() error {
	config := envget.ConfigInit()
	if len(config.JWTsecret) < 32 {
		return fmt.Errorf("JWT_SECRET must be set and at least 32 characters long")
	}
	jwtSecret := []byte(config.JWTsecret)

	storage, err := database.Connect(hasDropFlag(), &config)
	if err != nil {
		fmt.Println("Error connecting Postgres:", err)
		return err
	}
	defer storage.Close()

	hs := handlers.NewStorageHandler(storage, jwtSecret)

	r := gin.Default()

	err = users.CreateAdmin(context.Background(), &config, storage)
	if err != nil {
		return err
	}

	r.GET("/", hello)

	usersRouter := r.Group("/users")
	{
		usersRouter.POST("/login", hs.Login)
	}

	adminRouter := r.Group("/admin", handlers.AuthMiddleware(jwtSecret), handlers.AdminMiddleware())
	{
		adminRouter.POST("/User", hs.AddUser)
		adminRouter.GET("/User", hs.GetAllUsers)
	}

	fmt.Println("http://localhost:8080/")
	if err := r.Run(":8080"); err != nil {
		return err
	}
	return nil
}
