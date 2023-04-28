package main

import (
	"log"
	"net/http"
)

// DBConnection -- The database connection
var DBConnection string

// SessionTimeout -- Users session timeout
var SessionTimeout int

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	/*if len(os.Args) != 3 {
		log.Println("Usage__: " + os.Args[0] + " [full path to configfile] [environment (prod|dev|test)]")
		log.Println("Example: " + os.Args[0] + " /etc/grooming.json prod")
	} else {
		cFile := os.Args[1]
		cEnv := strings.ToLower(os.Args[2])
		if len(cFile) > 5 && len(cEnv) > 2 {
			f, err := os.Open(cFile)
			if err != nil {
				panic(err)
			}
			defer f.Close()
			config, err := gonfig.FromJson(f)
			if err != nil {
				panic(err)
			}
			//DBConnection, _ = config.GetString(cEnv+"/DBConnection", "sqlite3msr.db")
			DBConnection = "/Users/user/work/go/src/msr_members/sqlite3msr.db"
			port, _ := config.GetString(cEnv+"/Port", "8084")
			port = ":" + port
			SessionTimeout, _ = config.GetInt(cEnv+"/SessionTimeout", 16)*/
	port := ":8084"
	SessionTimeout = 10
	DBConnection = "./sqlite3msr.db"
	log.Printf("Port=%s, SessionTimeout=%d, DB=%s\n", port, SessionTimeout, DBConnection)
	router := NewRouter()
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))

	log.Println("Starting, listening to port " + port)
	log.Fatal(http.ListenAndServe(port, router))

	/*} else {
			log.Println("Config issue, Usage__: " + os.Args[0] + " [full path to configfile] [environment (prod|dev|test)]")
		}
	}*/
}
