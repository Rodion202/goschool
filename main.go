package main

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func hello(w http.ResponseWriter,r *http.Request){
	w.Write([]byte("Nothing interesting here yet but i added chi"))
}

func ConnToPostgres() (*pgxpool.Pool, error){
	user:=os.Getenv("POSTGRES_USER")
	pass:=os.Getenv("POSTGRES_PASSWORD")
	name:=os.Getenv("POSTGRES_DB")
	host:=os.Getenv("POSTGRES_HOST")

	if host==""{
		host="localhost:5432"
	}

	if user== "" || pass== "" || name == "" {
		return nil, fmt.Errorf("POSTGRES_USER,POSTGRES_PASSWORD and POSTGRES_DB must be set")
	}

	u:=url.URL{
		Scheme: "postgres",
		User: url.UserPassword(user,pass),
		Host: host,
		Path: name,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return pgxpool.New(ctx, u.String())
}

func main(){
	if err:=godotenv.Load();err!=nil{
		fmt.Println("no .env file found")
	}

	pool,err:=ConnToPostgres()
	if err!=nil{
		fmt.Println(err)
	} else {
		fmt.Println("Postgres is working!")
	}
	defer pool.Close()

	r:=chi.NewRouter()

	r.Get("/",hello)
	
	fmt.Println("http://localhost:8080/")
	if err:=http.ListenAndServe(":8080",r);err!=nil{
		fmt.Println("Oops:",err)
	}
}
