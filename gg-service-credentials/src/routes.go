package main

import "fmt"

func (s *Server) routes() {
	//Register endpoints for this microservice
	s.router.HandleFunc("/login", s.handleLogin()).Methods("GET")

	fmt.Println("Server routes setup successful...")
}
