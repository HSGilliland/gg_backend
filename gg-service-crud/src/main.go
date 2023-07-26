package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/lib/pq"
)

var dbconf Config

func init() {
	initConfiguration()
}

func initConfiguration() {
	dbconf = databaseConfiguration()
	fmt.Printf("CRUD Configuration Loaded --> %v:%v\n", dbconf.PostgresHost, dbconf.PostgresPort)
}

// Generate the DB Config
func databaseConfiguration() Config {
	conf := Config{
		UserName:     os.Getenv("POSTGRES_USER"),
		Password:     os.Getenv("POSTGRES_PASSWORD"),
		DatabaseName: os.Getenv("DATABASE"),
		Port:         os.Getenv("PORT"),
		PostgresHost: os.Getenv("POSTGRES_HOST"),
		PostgresPort: os.Getenv("POSTGRES_PORT"),
	}
	return conf
}

// THIS IS IT YO!
func main() {
	tmpConn, err := openDatabase(dbconf.PostgresHost, dbconf.PostgresPort, dbconf.UserName, dbconf.Password, dbconf.DatabaseName)
	if err != nil {
		log.Fatal(err)
	}

	server := Server{
		dbAccess: tmpConn,
		router:   mux.NewRouter(),
	}

	server.routes()
	handler := removeTrailingSlash(server.router)

	fmt.Printf("Starting server on port -> %v\n", dbconf.Port)
	log.Fatal(http.ListenAndServe(":"+dbconf.Port, handler))
}

func openDatabase(host, port, userName, password, databaseName string) (*sql.DB, error) {
	tmpDB, err := sql.Open("postgres", "user="+userName+" password="+password+" host="+host+" port="+port+" dbname="+databaseName+" sslmode=disable")
	if err != nil {
		fmt.Println("DB Crashed on open")
		return tmpDB, err
	}

	//ping database utnil it comes online, or fails
	err = tmpDB.Ping()
	for retry := 0; err != nil && retry < 20; err = tmpDB.Ping() {
		if err, ok := err.(*pq.Error); ok {
			fmt.Println("pq Error: ", err.Code.Name())
		}
		//fmt.Println("user=" + userName + " password=" + password + " host=" + host + " port=" + port + " dbname=" + databaseName + " sslmode=disable")
		fmt.Println("Sleeping till connection opens")
		time.Sleep(1 * time.Second)
		retry++
	}

	if err != nil {
		fmt.Println("Crashes waiting for connection")
		return tmpDB, err
	}
	fmt.Println("Opened database connection!!")
	return tmpDB, err
}

func removeTrailingSlash(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.URL.Path = strings.TrimSuffix(r.URL.Path, "/")
		next.ServeHTTP(w, r)
	})
}

func throwHttpError(w http.ResponseWriter, statusCode int, errMessage string) {
	w.WriteHeader(500)
	fmt.Fprintf(w, errMessage)
	fmt.Println(errMessage)
	return
}
