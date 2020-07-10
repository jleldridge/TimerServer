package main

import (
	"fmt"
	"log"
	"net/http"
	// "encoding/base64"
	"golang.org/x/crypto/bcrypt"
	dbActions "github.com/jleldridge/gotimerapp/db"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func main() {
	handleRequests()
}

func handleRequests() {
	r := mux.NewRouter()

	// unauthorized paths
	r.HandleFunc("/newUser", createUser).Methods("POST")

	// authorized paths
	r.HandleFunc("/start", basicAuth(startTimer)).Methods("POST")
	r.HandleFunc("/stop", basicAuth(stopTimer)).Methods("POST")
	r.HandleFunc("/updatePassword", basicAuth(updateUserPassword)).Methods("POST")

	fmt.Println("Server started!")
	log.Fatal(http.ListenAndServe(":10000", r))
}

func requirePassword () {
	
}

func startTimer (w http.ResponseWriter, r *http.Request) {
	project := "test"
	description := "test description"
	fmt.Println("Received start request")
	err := dbActions.StartTimer(project, description)
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}

	fmt.Fprintf(w, "Timer %s %s started.", project, description)
}

func stopTimer (w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received stop request")
	err := dbActions.StopTimer()
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