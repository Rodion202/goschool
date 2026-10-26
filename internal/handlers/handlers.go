package handlers

import (
	"goschool/internal/auth"
	"goschool/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type StorageHandler struct {
	s         models.Storage
	jwtSecret []byte
}

func comparePasswords(hashedPwd string, plainPwd string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPwd), []byte(plainPwd))
	return err == nil
}

type ErrorResp struct {
	Err string `json:"error"`
}

func NewStorageHandler(storage models.Storage, secret []byte) *StorageHandler {
	return &StorageHandler{s: storage, jwtSecret: secret}
}

func SendJsonError(c *gin.Context, message string, statusCode int) {
	c.JSON(statusCode, ErrorResp{Err: message})
}

func (hs *StorageHandler) Login(c *gin.Context) {
	var input models.LoginInput
	ctx := c.Request.Context()

	if err := c.ShouldBindJSON(&input); err != nil {
		SendJsonError(c, "invalid input", http.StatusBadRequest)
		return
	}

	user, err := hs.s.GetUserByEmail(ctx, input.Email)
	if err != nil || user == nil {
		SendJsonError(c, "invalid credentials", http.StatusUnauthorized)
		return
	}

	if !comparePasswords(user.Password, input.Password) {
		SendJsonError(c, "invalid credentials", http.StatusUnauthorized)
		return
	}

	jwtToken, err := GenerateJWT(hs.jwtSecret, user.Id, user.Role)
	if err != nil {
		SendJsonError(c, "could not generate token", http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"access_token": jwtToken,
		"token_type":   "bearer",
	})
}

func (hs *StorageHandler) GetAllUsers(c *gin.Context) {
	ctx := c.Request.Context()
	users, err := hs.s.GetAllUsers(ctx)
	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to GetAllUsers")
		return
	}
	usermessage := make([]models.UserMessage, len(users))
	for i, user := range users {
		usermessage[i] = models.UserMessage{Id: user.Id, Name: user.Name, Email: user.Email, Role: user.Role}
	}
	c.JSON(http.StatusOK, usermessage)
}

func (hs *StorageHandler) AddUser(c *gin.Context) {
	var user models.User
	ctx := c.Request.Context()
	if err := c.ShouldBindJSON(&user); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	hash, err := auth.HashPassword(user.Password)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	user.Password = hash

	err = hs.s.AddUser(ctx, &user)
	if err != nil {
		SendJsonError(c, "Внутренняя ошибка сервера. Попробуйте позже"+err.Error(), http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusCreated, models.UserMessage{
		Id:    user.Id,
		Name:  user.Name,
		Email: user.Email,
		Role:  user.Role,
	})
}
