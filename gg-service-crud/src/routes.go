package main

import "fmt"

func (s *Server) routes() {
	//Register and Login Endpoints
	s.router.HandleFunc("/login", s.handleLogin()).Methods("GET")
	s.router.HandleFunc("/register", s.handleRegisterPerson()).Methods("POST")

	fmt.Println("CRUD: Server routes setup successful...")
}
