package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"goschool/database"
	"goschool/envget"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
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
	var pool *pgxpool.Pool
	var err error
	if isArgumented() {
		pool,err=database.ConnectDrop(envget.GetPostgresReqs())
	} else {
	pool,err=database.Connect(envget.GetPostgresReqs())
}
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
