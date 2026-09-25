package main

import (
	"context"
	"fmt"
	"net/http"

	"goschool/database"
	"goschool/envget"
	"github.com/go-chi/chi/v5"
)

func hello(w http.ResponseWriter,r *http.Request){
	w.Write([]byte("Nothing interesting here yet but i added chi"))
}

func main(){
	pool,err:=database.Connect(envget.GetPostgresReqs())
	if err!=nil{
		fmt.Println("Error conecting Postgres:",err)
	}
	defer pool.Close()

	err=pool.Ping(context.Background())
	if err!=nil{
		fmt.Println("Ping error:",err)
	}

	r:=chi.NewRouter()

	r.Get("/",hello)
	
	fmt.Println("http://localhost:8080/")
	if err:=http.ListenAndServe(":8080",r);err!=nil{
		fmt.Println("Oops:",err)
	}
}
