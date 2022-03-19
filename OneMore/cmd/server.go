package main

import (
	"github.com/gorilla/mux"
	"net/http"
	s "user/pkg/service"
)
// The entry point

func main() {
	srv := s.NewService()
	router := mux.NewRouter()

	router.HandleFunc("/getUsers", srv.GetAllUsers)
	router.HandleFunc("/createUser", srv.CreateUser)
	router.HandleFunc("/{user_id:[0-9]+}/add_note/{tag_id:[0-9]+}", srv.AddNote)
	http.Handle("/", router)

	http.ListenAndServe("localhost:8080", nil)
}
