package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"
)

var credConfig Config

func init() {
	initConfiguration()
}

func initConfiguration() {
	credConfig = CredentialConfiguration()
	fmt.Println("Credentials Service configuration loaded....")
}

// CredentialConfiguration generates the config file used in this service
func CredentialConfiguration() Config {
	conf := Config{
		CredentialsPort: os.Getenv("CREDENTIALS_PORT"),
		CRUDHost:        os.Getenv("CRUD_Host"),
		CRUDPort:        os.Getenv("CRUD_Port"),
	}
	return conf
}

// This is it yo
func main() {
	server := Server{
		router: mux.NewRouter(),
	}

	//Create routes on this server
	server.routes()
	handler := removeTrailingSlash(server.router)

	fmt.Printf("Starting server on port -> %v\n", credConfig.CredentialsPort)
	log.Fatal(http.ListenAndServe(":"+credConfig.CredentialsPort, handler))
}

func buildLoginRequestURL(username string, password string) string {
	loginURL := "http://" + credConfig.CRUDHost + ":" + credConfig.CRUDPort + "/login" + "?username=" + username + "&password=" + password
	return loginURL
}

func throwHttpError(w http.ResponseWriter, statusCode int, errMessage string) {
	w.WriteHeader(500)
	fmt.Fprintf(w, errMessage)
	fmt.Println(errMessage)
	return
}

func removeTrailingSlash(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.URL.Path = strings.TrimSuffix(r.URL.Path, "/")
		next.ServeHTTP(w, r)
	})
}
