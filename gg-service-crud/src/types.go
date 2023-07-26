package main

import (
	"database/sql"

	"github.com/gorilla/mux"
)

// Config include the global settings
type Config struct {
	UserName     string
	Password     string
	DatabaseName string
	Port         string
	PostgresHost string
	PostgresPort string
}

// Server struct to manage server globally
type Server struct {
	dbAccess *sql.DB
	router   *mux.Router
}

type LoginResult struct {
	Username   string `json:"username"`
	UserExists string `json:"userexists"`
	PersonType string `json:"persontype"`
	PersonID   string `json:"personid"`
}

type RegisterPerson struct {
	Name       string `json:"name"`
	Surname    string `json:"surname"`
	Username   string `json:"username"`
	Password   string `json:"password"`
	PersonType string `json:"persontype"`
}

type RegisterPersonResponse struct {
	Status   string `json:"status"`
	Message  string `json:"message"`
	PersonId string `json:"personid"`
}
