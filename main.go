package main

import(
	"fmt"
	"net/http"
)

func hello(w http.ResponseWriter,r *http.Request){
	fmt.Fprintf(w,"Nothing interesting here yet")
}

func main(){
	mux:=http.NewServeMux()
	
	mux.HandleFunc("/",hello)
	
	fmt.Println("http://localhost:8080/")
	if err:=http.ListenAndServe(":8080",mux);err!=nil{
		fmt.Println("Oops:",err)
	}
}
