package main

import (
	"fmt"
	"log"
	"net/http"
	"encoding/json"

	"golang.org/x/crypto/bcrypt"
	dbActions "github.com/jleldridge/gotimerapp/db"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type TimerParams struct {
	UserID int `json:"userId"`
	Project string `json:"project"`
	Description string `json:"description"`
}

type UserIDResponse struct {
	UserID int `json:"userId"`
}

func main() {
	r := mux.NewRouter()

	// unauthorized paths
	r.HandleFunc("/newUser", createUser).Methods("POST")

	// authorized paths
	r.HandleFunc("/login", basicAuth(login)).Methods("GET")
	r.HandleFunc("/start", basicAuth(startTimer)).Methods("POST")
	r.HandleFunc("/stop", basicAuth(stopTimer)).Methods("POST")
	r.HandleFunc("/updatePassword", basicAuth(updateUserPassword)).Methods("POST")

	fmt.Println("Server started!")
	log.Fatal(http.ListenAndServe(":10000", r))
}

func login (w http.ResponseWriter, r *http.Request) {
	user, _, ok := r.BasicAuth()
	
	if !ok {
		w.WriteHeader(404)
		fmt.Fprintf(w, "User not found")
	}

	userId, err := dbActions.GetUserID(user)
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	js, err := json.Marshal(UserIDResponse{userId})
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }

  w.Header().Set("Content-Type", "application/json")
  w.Write(js)
}

func startTimer (w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var data TimerParams
	err := decoder.Decode(&data)
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}

	if data.UserID == 0 {
		w.WriteHeader(400)
		return
	}

	err = dbActions.StartTimer(data.UserID, data.Project, data.Description)
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}

	fmt.Fprintf(w, "Timer %s %s started.", data.Project, data.Description)
}

func stopTimer (w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var data TimerParams
	err := decoder.Decode(&data)
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}

	if data.UserID == 0 {
		w.WriteHeader(400)
		return
	}

	err = dbActions.StopTimer(data.UserID)
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}

	fmt.Fprintf(w, "Timer stopped.")
}

func createUser(w http.ResponseWriter, r *http.Request) {
	user, pass, ok := r.BasicAuth()
	if !ok {
		w.WriteHeader(400)
		return
	}
	
	fmt.Printf("Creating user %s with password %s\n", user, pass)

	err := dbActions.CreateUser(user, hashAndSalt([]byte(pass)))
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}

	fmt.Fprintf(w, "User %s created successfully!", user)
}

func updateUserPassword(w http.ResponseWriter, r *http.Request) {
}

func basicAuth(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, pass, ok := r.BasicAuth()
		hashedPass, err := dbActions.GetHashedPassword(user)
		if err != nil {
			fmt.Fprintf(w, err.Error())
			return
		}

		if !ok || !comparePasswords(hashedPass, pass) {
			w.Header().Set("WWW-Authenticate", `Basic realm="whatever"`)
			w.WriteHeader(401)
			w.Write([]byte("Unauthorized.\n"))
			return
		}

		handler(w, r)
	}
}


func hashAndSalt(pwd []byte) string {
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
			log.Println(err)
	}
	return string(hash)
}

func comparePasswords(hashedPwd string, plainPwd string) bool {
	byteHash := []byte(hashedPwd)
	bytePwd := []byte(plainPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, bytePwd)
	if err != nil {
			log.Println(err)
			return false
	}
	
	return true
}