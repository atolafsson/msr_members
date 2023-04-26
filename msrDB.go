package main

import (
	"database/sql"
	"log"
)

var DB *sql.DB
var dbOpen = false

const (
	sqlMembers = `SELECT Id,Name, Email, NickName, Prospect, Address, Zip, City, Notes FROM members`
)

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

// OpenDB : Open the Grooming database
func OpenDB() *sql.DB {
	sqlInfo := DBConnection
	db, err := sql.Open("sqlite3", sqlInfo)
	checkErr(err)
	err = db.Ping()
	checkErr(err)
	log.Println("DB is now Open")
	dbOpen = true
	return db
}

func GetMembers() MemberList {
	var lMembers MemberList
	if !dbOpen {
		DB = OpenDB()
	}
	log.Println("DB Query=" + sqlMembers)
	rows, err := DB.Query(sqlMembers)
	checkErr(err)
	defer rows.Close()
	log.Println("populate the data")
	for rows.Next() {
		var r Member
		err = rows.Scan(&r.id, &r.name, &r.email, &r.nickName, &r.prospect, &r.address, &r.zip, &r.city, &r.notes)
		checkErr(err)
		lMembers = append(lMembers, r)
	}
	return lMembers
}
