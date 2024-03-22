package main

import (
	"log"
	"net/http"
)

/***************************************************
To build a docker image
 $ docker volume create msr-db
 $ docker build -t msr-members .
 $ docker run -dp 0.0.0.0:8084:8084 --mount type=volume,src=msr-db,target=/etc/msr msr-members
****************************************************/

// DBConnection -- The database connection
var DBConnection string
var DBDir string

// SessionTimeout -- Users session timeout
var SessionTimeout int

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	port := ":8084"
	SessionTimeout = 10
	DBDir = "/etc/msr"
	DBConnection = DBDir + "/members.db"
	log.Printf("Port=%s, SessionTimeout=%d, DB=%s\n", port, SessionTimeout, DBConnection)
	router := NewRouter()
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))

	log.Println("Starting, listening to port " + port)
	log.Fatal(http.ListenAndServe(port, router))
}
