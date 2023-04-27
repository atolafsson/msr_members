package main

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB
var dbOpen = false

const (
	sqlMembers    = `SELECT Id, Name, Email, NickName, Prospect, Address, Zip, City, Notes FROM members`
	sqlMembersS   = `SELECT Name, Email, NickName FROM members`
	sqlGetMember  = `SELECT Id, Name, Email, NickName, Prospect, Address, Zip, City, Notes FROM members WHERE id = ?`
	sqlMemberSave = `UPDATE members set Name=?, Email=?, NickName=?, Prospect=?, Address=?, Zip=?, City=?, Notes=? where id=?`
	sqlMemberAdd  = `INSERT into members (Name, Email, NickName, Prospect, Address, Zip, City, Notes) values (?,?,?,?,?,?,?,?,?) RETURNING id`
)

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

// OpenDB : Open the Grooming database
func OpenDB() *sql.DB {
	sqlInfo := DBConnection
	log.Println("Open database sqlite3, " + sqlInfo)
	db, err := sql.Open("sqlite3", sqlInfo)
	checkErr(err)
	err = db.Ping()
	checkErr(err)
	log.Println("DB is now Open")
	dbOpen = true
	return db
}

func GetMember(mid int) Member {
	var r Member
	if !dbOpen {
		DB = OpenDB()
	}
	row := DB.QueryRow(sqlGetMember, mid)
	err := row.Scan(&r.id, &r.name, &r.email, &r.nickName, &r.prospect, &r.address, &r.zip, &r.city, &r.notes)
	checkErr(err)
	return r
}

func GetMembers() MemberList {
	var lMembers MemberList
	if !dbOpen {
		DB = OpenDB()
	}
	//ins, _ := DB.Prepare("INSERT  into members (Name, Email, NickName, Prospect, Notes, Address, Zip, City) VALUES (?, ?, ?, ?, ?, ?, ?, ?)")
	//ins.Exec("Halldor", "aatlason@gmail.com", "Axel", 0, "", "201 Countryside Ln", 38040, "Halls")
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

// UpdMember -- Update or add a single member
func UpdMember(mb Member) int {
	id := mb.id
	if !dbOpen {
		DB = OpenDB()
	}
	if id > 0 {
		_, err := DB.Exec(sqlMemberSave, mb.name, mb.email, mb.nickName, mb.prospect, mb.address, mb.zip, mb.city, mb.notes, mb.id)
		checkErr(err)
	} else {
		err := DB.QueryRow(sqlMemberAdd, mb.name, mb.email, mb.nickName, mb.prospect, mb.address, mb.zip, mb.city, mb.notes).Scan(&id)
		checkErr(err)
	}
	return id
}

func GetMembersS() MemberSList {
	var lMembers MemberSList
	if !dbOpen {
		DB = OpenDB()
	}
	log.Println("DB Query=" + sqlMembersS)
	rows, err := DB.Query(sqlMembersS)
	checkErr(err)
	defer rows.Close()
	log.Println("populate the data")
	for rows.Next() {
		var r MemberS
		err = rows.Scan(&r.name, &r.email, &r.nickName)
		checkErr(err)
		lMembers = append(lMembers, r)
	}
	return lMembers
}
