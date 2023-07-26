package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (s *Server) handleLogin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username := r.URL.Query().Get("username")
		password := r.URL.Query().Get("password")

		//call function from PSql into variables
		var validuser, persontype, personid string
		var querystring string
		querystring = "SELECT * FROM public.login('" + username + "','" + password + "')"
		err := s.dbAccess.QueryRow(querystring).Scan(&validuser, &persontype, &personid)

		//if queryrow returns error, provide error to caller and exit
		if err != nil {
			throwHttpError(w, 500, "/login -> Unable to query database...")
			return
		}

		loginResult := LoginResult{}
		loginResult.UserExists = validuser
		loginResult.Username = username
		loginResult.PersonType = persontype
		loginResult.PersonID = personid

		js, jserr := json.Marshal(loginResult)

		//if Json.marshal returns error, provide error to caller and exit
		if jserr != nil {
			throwHttpError(w, 500, "/login -> Unable to create JSON from DB result.")
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write(js)
	}
}

func (s *Server) handleRegisterPerson() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		regPerson := RegisterPerson{}
		err := json.NewDecoder(r.Body).Decode(&regPerson)
		if err != nil {
			throwHttpError(w, 500, "/register -> Bad Json Provided: "+err.Error())
			return
		}

		//Check if username already exists?
		var userNameExists string
		queryString := "SELECT * FROM public.registerverifyuser('" + regPerson.Username + "')"
		err = s.dbAccess.QueryRow(queryString).Scan(&userNameExists)
		//Check if queryrow returns error
		if err != nil {
			throwHttpError(w, 500, "/register -> Unable to query Database(RegisterVerifyUser)")
			return
		}

		fmt.Println("Username exists -> " + userNameExists)

		//Respond to user that the current username already exists
		if userNameExists == "true" {
			regPersonResponse := RegisterPersonResponse{}
			regPersonResponse.Status = "ERROR"
			regPersonResponse.Message = "This username is already in use."
			regPersonResponse.PersonId = ""

			js, jserr := json.Marshal(regPersonResponse)
			if jserr != nil {
				throwHttpError(w, 500, "/register -> Unable to create JSON from RegisterPersonResponse - Verify")
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(200)
			w.Write(js)
		}

		//If username does not exist, continue to create a new person
		var personCreated bool
		var username, personId string
		queryString = "SELECT * FROM public.registerperson('" + regPerson.Name + "','" + regPerson.Surname + "','" + regPerson.Username + "','" + regPerson.Password + "','" + regPerson.PersonType + "')"

		err = s.dbAccess.QueryRow(queryString).Scan(&personCreated, &username, &personId)

		if err != nil {
			throwHttpError(w, 500, "/login -> Unable to process DB Function (RegisterPerson). Error: "+err.Error())
			return
		}

		regPersonResponse := RegisterPersonResponse{}
		if personCreated {
			regPersonResponse.Status = "SUCCESS"
			regPersonResponse.Message = "User created successfully."
			regPersonResponse.PersonId = personId
		} else {
			regPersonResponse.Status = "ERROR"
			regPersonResponse.Message = "Unable to create new user."
			regPersonResponse.PersonId = ""
		}

		js, jserr := json.Marshal(regPersonResponse)

		if jserr != nil {
			throwHttpError(w, 500, "/register -> Unable to Marshal Json (After Creation)")
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write(js)

	}
}
