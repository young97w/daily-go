package main

import (
	"fmt"
	"geektime/02type_server/server_context"
	"net/http"
)

func main() {
	//home
	server := server_context.NewServer("myServer")
	//server.Route("/", homeHandler)
	//server.Route("/user", userHandler)
	server.Route(http.MethodPost, "/signup", server_context.SignUp)
	err := server.Start(":8080")
	if err != nil {
		panic(err)
	}
}

func userHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcom to user page!")
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcom to home page!")
}
