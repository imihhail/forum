package main

import (
	"db/functions"
	"fmt"
	"net/http"
)

func main() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/home", functions.Homepage)
	http.HandleFunc("/register", functions.Register)
	http.HandleFunc("/comment", functions.CommentsPage)
	http.HandleFunc("/login", functions.Login)
	http.HandleFunc("/logout", functions.Logout)
	http.HandleFunc("/", functions.Handle404)
	//http.HandleFunc("/", functions.Homepage)
	fmt.Println("Server is running at: http://localhost:8080/home")
	fmt.Println()
	http.ListenAndServe("0.0.0.0:8080", nil)
}
