package main

import (
	"fmt"
	"log"
	"net/http"
	dbActions "github.com/jleldridge/gotimerapp/db"

	_ "github.com/lib/pq"
)

func main() {
	handleRequests()
}

func handleRequests() {
	http.HandleFunc("/start", startTimer)
	http.HandleFunc("/stop", stopTimer)
	http.HandleFunc("/entries", queryRows)
	fmt.Println("Server started!")
	log.Fatal(http.ListenAndServe(":10000", nil))
}

func startTimer(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received start request")
	err := dbActions.StartTimer("test", "test description")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}
}

func stopTimer(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received stop request")
	err := dbActions.StopTimer()
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}
}

func queryRows(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Querying rows")
	res, err := dbActions.QueryRows()
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}

	resString := ""
	for res.Next() {
		var id int64
		var start string
		var stop string
		var desc string
		var project string
		res.Scan(&id, &start, &stop, &desc, &project)
		resString += fmt.Sprintf("%d %s %s %s %s\n", id, start, stop, desc, project)
	}
	fmt.Fprintf(w, resString)
}