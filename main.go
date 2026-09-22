package main

import(
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func hello(w http.ResponseWriter,r *http.Request){
	w.Write([]byte("Nothing interesting here yet but i added chi"))
}

func main(){
	r:=chi.NewRouter()

	r.Get("/",hello)
	
	fmt.Println("http://localhost:8080/")
	if err:=http.ListenAndServe(":8080",r);err!=nil{
		fmt.Println("Oops:",err)
	}
}
