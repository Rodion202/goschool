package handlers

import (
	"encoding/json"
	"goschool/models"
	"net/http"
	//"github.com/go-chi/chi/v5"
)

type StorageHandler struct{
	s models.Storage
}

type ErrorResp struct{
	Err string `json:"error"`
}

func NewStorageHandler (storage models.Storage) *StorageHandler{
	return &StorageHandler{s: storage}
}

func SendJsonError (w http.ResponseWriter, message string, statusCode int){
	w.Header().Set("Content-Type","application/`json:")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(ErrorResp{Err: message})
}

func (hs StorageHandler) GetAllUsers (w http.ResponseWriter, r *http.Request){
	users,err:=hs.s.GetAllUsers()
	if err!=nil{
		http.Error(w,"Failed to GetAllUsers",http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type","application/json")
	json.NewEncoder(w).Encode(users)
}

func (hs StorageHandler) AddUser (w http.ResponseWriter,r *http.Request){
	var user models.User
	if err:=json.NewDecoder(r.Body).Decode(&user);err!=nil{
		http.Error(w,err.Error(),http.StatusBadRequest)
		return
	}
	err:=hs.s.AddUser(user)
	if err!=nil{
		SendJsonError(w, "Внутренняя ошибка сервера. Попробуйте позже", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}


