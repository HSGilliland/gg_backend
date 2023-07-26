package main

import (
	"github.com/gorilla/mux"
)

// Server struct to manage server globally
type Server struct {
	router *mux.Router
}

// Config is the struct for the global connection settings
type Config struct {
	CredentialsPort string
	CRUDHost        string
	CRUDPort        string
}

// LoginUserRequest the data being sent to the crud to verify login creds
type LoginUserRequest struct {
	Username string
	Password string
}

// LoginUserResponse the data returned from the crud after login request
type LoginUserResponse struct {
	Username   string `json:"username"`
	UserExists string `json:"userexists"`
	PersonType string `json:"persontype"`
	PersonID   string `json:"personid"`
}
