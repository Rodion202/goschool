package main

import (
	"fmt"
	"net/http"
	"os"

	"goschool/database"
	"goschool/handlers"
	//"goschool/envget"

	"github.com/go-chi/chi/v5"
)

func hello(w http.ResponseWriter,r *http.Request){
	w.Write([]byte("Nothing interesting here yet but i added chi"))
}

func isArgumented()bool{
	if len(os.Args)<2{
		return false
	} else {
		if os.Args[1]=="drop"{
			return true
		}
	}
	return false
}

func main(){
	storage,err:=database.Connect(isArgumented())
	if err!=nil{
		fmt.Println("Error conecting Postgres:",err)
	}
	defer storage.Close()

	hs:=handlers.NewStorageHandler(storage)

	r:=chi.NewRouter()

	r.Get("/",hello)

	r.Route("/users",func(r chi.Router) {
	r.Get("/",hs.GetAllUsers)
	r.Post("/",hs.AddUser)
})
	
	fmt.Println("http://localhost:8080/")
	if err:=http.ListenAndServe(":8080",r);err!=nil{
		fmt.Println("Oops:",err)
	}
}
