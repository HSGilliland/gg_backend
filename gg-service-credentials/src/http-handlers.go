package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func (s *Server) handleLogin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username := r.URL.Query().Get("username")
		password := r.URL.Query().Get("password")

		fmt.Println("Login Called <---> Username -> " + username + "Password -> " + password)

		//TODO: This should return a JSON Payload, not just text
		//TODO: Create a list of return codes, should not just return a generic 400
		//check if a username was supplied
		if username == "" {
			throwHttpError(w, 400, "/login -> No username was supplied")
			return
		}

		//Checks if password was supplied
		if password == "" {
			throwHttpError(w, 400, "/login -> No password was supplied")
			return
		}

		loginUrl := buildLoginRequestURL(username, password)

		//Call Login Service
		req, err := http.Get(loginUrl)

		if err != nil {
			fmt.Println("1")
			throwHttpError(w, 500, err.Error())
			return
		}

		defer req.Body.Close()
		if req.StatusCode == 500 {
			bodyBytes, err := ioutil.ReadAll(req.Body)
			if err != nil {
				log.Fatal(err)
			}
			bodyString := string(bodyBytes)
			throwHttpError(w, 500, "/login -> Request to database could not be completed --> "+bodyString)
			return
		}

		var loginResponse LoginUserResponse

		decoder := json.NewDecoder(req.Body)

		//Check if there was an error while trying to decode the JSON response
		err = decoder.Decode(&loginResponse)
		if err != nil {
			fmt.Println("2")
			throwHttpError(w, 500, err.Error())
			return
		}

		js, jserr := json.Marshal(loginResponse)
		if jserr != nil {
			throwHttpError(w, 500, "/login -> Unable to create JSON from Login Result.")
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write(js)
	}
}
