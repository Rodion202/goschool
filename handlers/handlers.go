package handlers

import (
	"goschool/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type StorageHandler struct {
	s models.Storage
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), err
}

func comparePasswords(hashedPwd string, plainPwd string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPwd), []byte(plainPwd))
	return err == nil
}

type ErrorResp struct {
	Err string `json:"error"`
}

func NewStorageHandler(storage models.Storage) *StorageHandler {
	return &StorageHandler{s: storage}
}

func SendJsonError(c *gin.Context, message string, statusCode int) {
	c.JSON(statusCode, ErrorResp{Err: message})
}

func (hs *StorageHandler) Login(c *gin.Context) {
	var input models.LoginInput

	if err := c.ShouldBindJSON(&input); err != nil {
		SendJsonError(c, "invalid input", http.StatusBadRequest)
		return
	}

	user, err := hs.s.GetUserByEmail(input.Email)
	if err != nil {
		SendJsonError(c, "invalid credentials", http.StatusUnauthorized)
		return
	}

	if !comparePasswords(user.Password, input.Password) {
		SendJsonError(c, "invalid credentials", http.StatusUnauthorized)
		return
	}

	jwtToken, err := GenerateJWT(user.Id, user.Role)
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
	users, err := hs.s.GetAllUsers()
	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to GetAllUsers")
		return
	}
	usermessage := make([]models.UserMessage, len(users))
	for i, user := range users {
		usermessage[i] = models.UserMessage{Id: user.Id, Name: user.Name, Email: user.Email}
	}
	c.JSON(http.StatusOK, usermessage)
}

func (hs *StorageHandler) AddUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	hash, err := hashPassword(user.Password)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	user.Password = hash

	err = hs.s.AddUser(&user)
	if err != nil {
		SendJsonError(c, "Внутренняя ошибка сервера. Попробуйте позже", http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusCreated, models.UserMessage{
		Id:    user.Id,
		Name:  user.Name,
		Email: user.Email,
	})
}
